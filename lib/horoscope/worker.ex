defmodule Horoscope.Worker do
  use GenServer
  alias Horoscope.Repo

  # Public API
  def get, do: get(:calendar.iso_week_number())
  def get(sign) when is_binary(sign) do
    :calendar.iso_week_number()
    |> Tuple.insert_at(2, sign)
    |> get
  end
  def get(params) when tuple_size(params) in [2, 3] do
    GenServer.call(name, params)
  end

  # GenServer API
  def start_link do
    GenServer.start_link(__MODULE__, [], [name: name])
  end

  def handle_call({year, week}, _from, state) do
    results = horoscope(year, week) |> Repo.all
    {:reply, results, state}
  end

  def handle_call({year, week, sign}, _from, state) do
    result = horoscope(year, week, sign) |> Repo.one
    {:reply, result, state}
  end

  # Function API
  import Ecto.Query
  alias Horoscope.Model

  defp horoscope(year, week) do
    Model
    |> where([h], h.year == ^year and h.week == ^week)
    |> select([h], %{id: h.id, sign: h.sign, prediction: h.prediction})
  end

  defp horoscope(year, week, sign) do
    sign = sign |> to_string |> String.downcase |> String.capitalize

    horoscope(year, week)
    |> where([h], h.sign == ^sign)
  end

  # Config
  @name :horoscope_worker
  def name, do: @name
end