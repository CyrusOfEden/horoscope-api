defmodule Horoscope.Mixfile do
  use Mix.Project

  def project do
    [app: :horoscope,
     version: "0.0.1",
     elixir: "~> 1.0",
     build_embedded: Mix.env == :prod,
     start_permanent: Mix.env == :prod,
     deps: deps]
  end

  def application do
    [mod: {Horoscope, []},
     applications: [:logger, :quantum,
                    :cowboy, :plug,
                    :poolboy, :httpoison,
                    :postgrex, :ecto]]
  end

  defp deps do
    [{:cowboy, "~> 1.0"},
     {:plug, "~> 1.0"},
     {:poolboy, "~> 1.5"},
     {:httpoison, "~> 0.7"},
     {:poison, "~> 1.5"},
     {:floki, "~> 0.3"},
     {:ecto, "~> 1.0"},
     {:postgrex, "~> 0.9"},
     {:towel, "~> 0.2"},
     {:quantum, "~> 1.4"},
     {:churchill, github: "knrz/churchill"}]
  end
end
