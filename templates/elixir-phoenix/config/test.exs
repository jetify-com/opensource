import Config

# We don't run a server during test. If one is required,
# you can enable the server option below.
config :elixir_phoenix, ElixirPhoenixWeb.Endpoint,
  http: [ip: {127, 0, 0, 1}, port: 4002],
  secret_key_base: "bQFSwYiJNlm9t+feCW/Q2/iQR/5csXEzQnPKPKcr3B2diC5yoi1lBafJ4jJoKj6c",
  server: false

# In test we don't send emails.
config :elixir_phoenix, ElixirPhoenix.Mailer, adapter: Swoosh.Adapters.Test

# Disable swoosh api client as it is only required for production adapters.
config :swoosh, :api_client, false

# Print only warnings and errors during test
config :logger, level: :warning

# Initialize plugs at runtime for faster test compilation
config :phoenix, :plug_init_mode, :runtime

config :phoenix_live_view,
  # Enable helpful, but potentially expensive runtime checks
  enable_expensive_runtime_checks: true
