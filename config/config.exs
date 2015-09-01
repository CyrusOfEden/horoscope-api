use Mix.Config

config :quantum,
  timezone: :local,
  cron: ["0 0 * * 0": {"Horoscope.Scraper", :fetch}]

import_config "#{Mix.env}.exs"
