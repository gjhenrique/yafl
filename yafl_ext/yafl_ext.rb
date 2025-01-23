#!/usr/bin/env ruby

require "fileutils"
require "json"
require "cgi"

class Bookmark
  QUERY=<<SQL
    WITH RECURSIVE parents(id) AS (
     SELECT id FROM moz_bookmarks WHERE (type = 2 AND title = 'Work')
     UNION ALL
     SELECT mb.id FROM moz_bookmarks mb INNER JOIN parents ON mb.parent = parents.id
    )
    SELECT DISTINCT mp.url, mb.title FROM moz_bookmarks mb
     INNER JOIN moz_places mp ON mp.id = mb.fk
     INNER JOIN parents ON mb.id = parents.id
     WHERE mb.type = 1 AND mb.title IS NOT NULL
     ORDER BY mp.frecency DESC
SQL

  def list
    require "sqlite3"

    ff_dir = Dir["#{Dir.home}/.var/app/org.mozilla.firefox/.mozilla/firefox/*default-release"][0]

    if not ff_dir
      pattern = File.expand_path('~/.mozilla/firefox/*.default*/')
      ff_dir = Dir.glob(pattern).find { |path| File.directory?(path) }
    end

    FileUtils.cp("#{ff_dir}/places.sqlite", "/tmp/places.sqlite")
    begin
      FileUtils.cp("#{ff_dir}/places.sqlite-wal", "/tmp/places.sqlite-wal")
    rescue
    end

    array = []
    db = SQLite3::Database.new "/tmp/places.sqlite"
    db.execute QUERY do |row|
      array << "#{row[0]}\x1f#{row[1]}"
    end

    return array
  end
end

class BookmarkOpen < Bookmark
  def execute(input)
    `xdg-open "#{input}"`
  end
end

class BookmarkCopy < Bookmark
  def execute(input)
    `wl-copy "#{input}"`
  end
end

class Search
  def execute(input)
    mode = config.find { |c| input.start_with? c[:trigger] }
    exit 1 if mode.nil?

    input = CGI.escape(input.gsub("#{mode[:trigger]} ", ""))
    url = mode[:url].gsub("#{mode[:trigger]} ", "").gsub("%s", input)

    `xdg-open "#{url}"`
  end

  def list
    return config.sort_by { |c| c[:trigger] }.map { |c| "#{c[:trigger]} #{c[:name]}" }
  end

  private
  def parse_json(filename)
    JSON.parse(File.read("#{Dir.home}/#{filename}"), symbolize_names: true)
  end

  def config
    parse_json(".config/yafl/search.json").tap do |config|
      config.push(*parse_json("Life/work/search.json")); rescue nil
    end
  end
end

class Theme
  def list
    [
      "Turn theme on",
      "Turn theme off",
    ]
  end

  def execute(input)
    if input.end_with? "on"
      `switch_theme on`
    elsif input.end_with? "off"
      `switch_theme off`
    end
  end
end

class Sunset
  def list
    [
      "Turn sunset on",
      "Turn sunset off",
    ]
  end

  def execute(input)
    if input.end_with? "on"
      `wlsunset`
    elsif input.end_with? "off"
      `pkill wlsunset`
    end
  end
end

class Passage
  require 'find'

  def list
    root_path = "/home/guilherme/.passage/store"
    result = []

    Find.find(root_path) do |path|
      next if path == root_path || path.include?(".git") || path.include?(".age-recipients")

      relative_path = path.sub("#{File.dirname(root_path)}/", '').gsub(".age", "").gsub("store/", "")
      result << [relative_path] if File.file?(path) || Dir.empty?(path)
    end

    result
  end

  def execute(input)
    `wl-copy $(passage #{input})`
  end
end

entry = ARGV[0].split('_').map(&:capitalize).join
handler = Module.const_get(entry).new

if ARGV.length > 1
  input = ARGV[1]
  handler.execute(input)
else
  puts handler.list.join("\n")
end
