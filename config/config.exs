use Mix.Config

work = fn ->
  apply("Horoscope.Scraper", :changesets, [true])
  |> Stream.take(12)
  |> Stream.run
end

config :quantum,
  timezone: :local,
  cron: ["0 0 * * 0": work]

config :horoscope, Horoscope.Repo,
  adapter: Ecto.Adapters.Postgres,
  database: "horoscope",
  username: System.get_env("DB_USERNAME"),
  password: System.get_env("DB_PASSWORD"),
  hostname: "localhost",
  pool_size: (if Mix.env == :prod, do: 4, else: 2)
