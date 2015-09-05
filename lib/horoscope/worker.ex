defmodule Horoscope.Worker do
  use GenServer
  use Towel
  import Horoscope, only: [worker_pool: 0]

  alias Horoscope.Repo

  # Public API
  @default_opts %{encode: false}
  def call(opts \\ %{}) do
    Map.merge(@default_opts, opts)
    |> Map.put_new(:week, :calendar.iso_week_number())
    |> Map.put_new(:sign, nil)
    |> process
  end

  def process(%{sign: nil, week: {year, week}, encode: encode}) do
    fetch({year, week}, encode)
  end
  def process(%{sign: sign, week: {year, week}, encode: encode}) do
    fetch({year, week, sign}, encode)
  end

  # Memoize
  def fetch(params, encode) do
    data = :poolboy.transaction(worker_pool, &GenServer.call(&1, params))
    case encode do
      true  -> Poison.encode!(data)
      false -> data
    end
  end

  # GenServer API
  def start_link(_) do
    GenServer.start_link(__MODULE__, [], [])
  end

  def handle_call({year, week}, _from, state) do
    horoscope(year, week)
    |> Maybe.wrap
    |> fmap(&Repo.all/1)
    |> handle_response(state)
  end

  def handle_call({year, week, sign}, _from, state) do
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

  @signs ~w[aquarius aries cancer capricorn
            gemini leo libra pisces
            sagittarius scorpio taurus virgo]
  defp horoscope(_, _, sign) when not sign in @signs, do: nothing
  defp horoscope(year, week, sign) do
    horoscope(year, week)
    |> where([h], h.sign == ^sign)
  end

  defp normalize_sign(sign) do
    sign |> to_string |> String.downcase
  end
end