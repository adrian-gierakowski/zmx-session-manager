{
  description = "zmx session manager";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";
    flake-compat = {
      url = "github:NixOS/flake-compat";
      flake = false;
    };
  };

  outputs = { self, nixpkgs, utils, ... }:
    let
      overlay = final: prev: {
        zsm = final.buildGoModule {
          pname = "zsm";
          version = "0.3.2";
          src = ./.;

          vendorHash = "sha256-ObIKdIvZ0TLBGC2C25MH1WnaqS+B0YhJmfUzCMrlgG8=";

          ldflags = [
            "-s"
            "-w"
            "-X main.version=0.3.2"
            "-X main.commit=none"
            "-X main.date=unknown"
          ];

          subPackages = [ "." ];

          postInstall = ''
            mv $out/bin/zmx-session-manager $out/bin/zsm
          '';

          meta = with final.lib; {
            description = "zmx session manager";
            homepage = "https://github.com/mdsakalu/zmx-session-manager";
            license = licenses.mit;
            mainProgram = "zsm";
          };
        };
      };
    in
    {
      overlays.default = overlay;
    } // utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [ overlay ];
        };
      in
      {
        packages.default = pkgs.zsm;

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gofumpt
            golangci-lint
          ];
        };
      }
    );
}
