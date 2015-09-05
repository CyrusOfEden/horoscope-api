use Mix.Config

config :quantum,
  timezone: :local,
  cron: ["0 0 * * 0": {"Horoscope.Scraper", :seed}]

import_config "#{Mix.env}.exs"
