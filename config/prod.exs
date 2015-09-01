use Mix.Config

config :data, Data.Repo,
  adapter: Ecto.Adapters.Postgres,
  database: "horoscope",
  username: "postgres",
  password: "",
  hostname: "localhost",
  pool_size: 6
