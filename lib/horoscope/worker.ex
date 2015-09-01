defmodule Horoscope.Worker do
  use GenServer
  use Towel

  import DefMemo
  require Logger

  alias Horoscope.Repo

  # Public API
  def call, do: fetch(:calendar.iso_week_number())
  def call(sign) when is_binary(sign) do
    {year, week} = :calendar.iso_week_number()
    fetch({year, week, sign})
  end
  def call(params) when tuple_size(params) in [2, 3] do
    fetch(params)
  end

  defmemo fetch(params) do
    :poolboy.transaction(Horoscope.worker_pool, &GenServer.call(&1, params))
  end

  # GenServer API
  def start_link(_) do
    GenServer.start_link(__MODULE__, [], [])
  end

  def handle_call({year, week}, _from, state) do
    Logger.debug("Getting horoscopes for week #{week} of #{year}")
    horoscope(year, week)
    |> Maybe.wrap
    |> fmap(&Repo.all/1)
    |> handle_response(state)
  end

  def handle_call({year, week, sign}, _from, state) do
    Logger.debug("Getting #{sign} horoscope for week #{week} of #{year}")
    horoscope(year, week, normalize_sign(sign))
    |> Maybe.wrap
    |> fmap(&Repo.one/1)
    |> handle_response(state)
  end

  defp handle_response(response, state) do
    case response do
      {:just, result} -> {:reply, result, state}
      :nothing        -> {:reply, nil, state}
    end
  end

  # Function API
  import Ecto.Query
  alias Horoscope.Model

  defp horoscope(year, week) do
    Model
    |> where([h], h.year == ^year and h.week == ^week)
    |> select([h], %{sign: h.sign, prediction: h.prediction})
  end

  @signs ~w[Aquarius Aries Cancer Capricorn
            Gemini Leo Libra Pisces
            Sagittarius Scorpio Taurus Virgo]
  defp horoscope(_, _, sign) when not sign in @signs, do: nothing
  defp horoscope(year, week, sign) do
    horoscope(year, week)
    |> where([h], h.sign == ^sign)
  end

  defp normalize_sign(sign) do
    sign |> to_string |> String.downcase |> String.capitalize
  end
end