{
  description = "FS Tracer Devshell";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
        };
      in
      with pkgs;
      {
        devShells.default = mkShell {
          buildInputs = [
            helmsman
            (wrapHelm kubernetes-helm {
              plugins = with kubernetes-helmPlugins; [
                helm-diff
              ];
            })
            bazel
            buildozer
            go
            gopls
            mockgen
            supabase-cli
          ];
        };
      }
    );
}
