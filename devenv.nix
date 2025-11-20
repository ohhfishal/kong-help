{ pkgs, lib, config, ... }@inputs:
let 
  pkgs-unstable = import inputs.nixpkgs-unstable { system = pkgs.stdenv.system; };
in
{
  languages.go = {
    enable = true;
    # package = pkgs-unstable.go;
  };

  git-hooks.hooks = {
    # Shell
    shellcheck.enable = true;

    # Golang
    govet.enable = true;
    gotest.enable = true;
    gofmt.enable = true;
    golangci-lint.enable = true;
  };
}
