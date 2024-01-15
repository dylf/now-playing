{
  description = "An API to get the currently playing song";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = [ "x86_64-linux" "aarch64-linux" "aarch64-darwin" "x86_64-darwin" ];

      perSystem = { config, self', inputs', pkgs, system, ... }:
        let
          name = "now-playing";
          version = "latest";
          vendorHash = null; # update whenever go.mod changes
        in
        {
          devShells = {
            default = pkgs.mkShell {
              inputsFrom = [ self'.packages.default ];
              packages = [ pkgs.postgresql_16 ];
              shellHook = ''
                export PGHOST=$HOME/postgres
                export PGDATA=$PGHOST/data
                export PGDATABASE=postgres
                export PGLOG=$PGHOST/postgres.log

                mkdir -p $PGHOST

                if [ ! -d $PGDATA ]; then
                  initdb --auth=trust --no-locale --encoding=UTF8
                fi

                if ! pg_ctl status
                then
                  pg_ctl start -l $PGLOG -o "--unix_socket_directories='$PGHOST'"
                fi
              '';
            };
          };

          packages = {
            default = pkgs.buildGoModule {
              inherit name vendorHash;
              src = ./.;
            };

            docker = pkgs.dockerTools.buildImage {
              inherit name;
              tag = version;
              config = {
                Cmd = "${self'.packages.default}/bin/${name}";
              };
            };
          };
        };
    };
}
