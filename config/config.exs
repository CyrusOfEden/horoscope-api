use Mix.Config

work = fn ->
  apply("Horoscope.Scraper", :stream, [true])
  |> Stream.take(12)
  |> Stream.run
end

config :quantum,
  timezone: :local,
  cron: ["0 0 * * 0": work]

import_config "#{Mix.env}.exs"
