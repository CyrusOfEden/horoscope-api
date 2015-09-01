use Mix.Config

config :horoscope, Horoscope.Repo,
  adapter: Ecto.Adapters.Postgres,
  database: "horoscope",
  username: "postgres",
  password: "",
  hostname: "localhost",
  pool_size: 8
