use Mix.Config

work = fn ->
  apply("Horoscope.Scraper", :seed, [])
  |> Stream.run
end

config :quantum,
  timezone: :local,
  cron: ["0 0 * * 0": work]

import_config "#{Mix.env}.exs"
