{
  inputs.nixpkgs.url = "github:nixos/nixpkgs/release-25.05";

  outputs =
    { nixpkgs, ... }:
    let
      systems = nixpkgs.lib.systems.flakeExposed;
    in
    {
      devShells = nixpkgs.lib.genAttrs systems (
        system:
        let
          pkgs = import nixpkgs { inherit system; };
        in
        {
          default = pkgs.mkShell {
            packages = with pkgs; [
              go
              gotools
              gopls
              just
              sqlc
              goose
            ];
          };
        }
      );
    };
}
