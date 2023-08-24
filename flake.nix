{
  description = "Discord bot for Pepsy's Gang";

  # Nixpkgs / NixOS version to use.
  inputs.nixpkgs.url = "nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs }:
    let

      # System types to support.
      supportedSystems =
        [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];

      # Helper function to generate an attrset '{ x86_64-linux = f "x86_64-linux"; ... }'.
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

      # Nixpkgs instantiated for supported system types.
      nixpkgsFor = forAllSystems (system:
        import nixpkgs {
          inherit system;
          overlays = [ self.overlays.default ];
        });
    in
    {

      # A Nixpkgs overlay.
      overlays.default = final: prev: {
        pepsy-discord-bot = with final;
          buildGoModule {

            pname = "pepsy-discord-bot";
            version = "v1.0";
            src = ./.;
            vendorSha256 = "sha256-t9unMXkACb+nypLC8asVRaSjXXaFNP51ULuBldDurEk=";

            meta = with lib; {
              description = "Discord bot for Pepsy's Gang";
              # description = "Discord bot for Pepsy";
              homepage = "https://github.com/pinpox/pepsy-discord-bot";
              license = licenses.gpl3;
              maintainers = with maintainers; [ pinpox ];
            };
          };
      };

      # Package
      packages = forAllSystems (system: {
        inherit (nixpkgsFor.${system}) pepsy-discord-bot;
        default = self.packages.${system}.pepsy-discord-bot;
      });

      # Nixos module
      nixosModules.pepsy-discord-bot = { pkgs, lib, config, ... }:
        with lib;
        let cfg = config.services.pepsy-discord-bot;
        in {

          options.services.pepsy-discord-bot = {
            # Optios for configuration
            enable = mkEnableOption "Pepsy Bot";
            discordTokenFile = mkOption {
              type = types.path;
              default = null;
              example = "/path/to/token";
              description = ''File containing Discord token with appropiate permissions'';
            };
          };

          config = mkIf cfg.enable {

            nixpkgs.overlays = [ self.overlays.default ];

            # Service
            systemd.services.pepsy-discord-bot = {

              # environment.DISCORD_TOKEN = "$(cat %d/discord_token)";

              wantedBy = [ "multi-user.target" ];
              after = [ "network.target" ];
              description = "Start the pepsy-bot";

              script = ''
                export DISCORD_TOKEN="$(cat ''${CREDENTIALS_DIRECTORY}/discord_token)"
                ${pkgs.pepsy-discord-bot}/bin/pepsy-discord-bot
              '';

              serviceConfig = {
                LoadCredential = [ "discord_token:${cfg.discordTokenFile}" ];

                WorkingDirectory = "${pkgs.pepsy-discord-bot}/bin";
                # User = "pepsy-discord-bot";
                # Group = "pepsy-discord-bot";
                DynamicUser = true;
                # ExecStart = "${pkgs.pepsy-discord-bot}/bin/pepsy-discord-bot";
              };
            };

          };
        };
    };
}
