{ lib, buildRubyGem, bundlerEnv, ruby }:

let 
  gems = bundlerEnv {
    inherit ruby;

    name = "yafl_ext";
    gemdir = ./.;
  };
in buildRubyGem rec {
  inherit ruby;

  name = "yafl_ext";

  gemName = "yafl_ext";
  pname = "yafl_ext";

  buildInputs = [gems];
}
