defmodule Horoscope.Scraper do
  use Towel

  alias Horoscope.Repo
  alias Horoscope.Model

  @base "http://www.theonion.com/features/horoscope"
  def fetch do
    case request(@base) do
      {:ok, page} ->
        page
        |> Floki.find(".reading-list-item")
        |> Floki.attribute("data-absolute-url")
        |> Stream.map(&get/1)
        |> Stream.flat_map(fn {:ok, {horoscopes, date}} ->
          Enum.map(horoscopes, &params(&1, date))
        end)
      {:error, r} ->
        error(r)
    end
  end

  def seed(count \\ 12, save \\ false)
  def seed(:all, save) do
    fetch
    |> Stream.map(&Model.changeset(%Model{}, &1))
    |> persist(save)
  end
  def seed(count, save) when is_integer(count) and count > 0 do
    seed(:all)
    |> Stream.take(count)
    |> persist(save)
  end

  defp persist(horoscopes, false), do: horoscopes
  defp persist(horoscopes, true) do
    horoscopes |> Stream.map(&Repo.insert/1)
  end

  def get(url) do
    case request(url) do
      {:ok, page} ->
        horoscopes =
          page
          |> Floki.find(".astro")
          |> Floki.find(".large-thing")
          |> Enum.map(fn elem -> elem |> Floki.text |> normalize end)
        date =
          page
          |> Floki.find(".content-published")
          |> Floki.text
          |> String.strip
          |> parse_date

        ok({horoscopes, date})
      {:error, r} ->
        error(r)
    end
  end

  defp params({sign, prediction}, date) do
    {year, week} = :calendar.iso_week_number(date)
    %{
      sign: sign,
      prediction: prediction,
      date: date,
      week: week,
      year: year
    }
  end

  defp request(url) do
    HTTPoison.get(url)
    |> fmap(&Map.get(&1, :body))
    |> fmap(&Floki.parse/1)
  end

  defp normalize(horoscope) do
    [details, prediction] =
      horoscope
      |> String.strip
      |> String.split(~r/[\s]{2,}/)
      |> Enum.map(&String.strip/1)
    sign =
      details
      |> String.split(" | ")
      |> List.first
      |> String.strip

    {sign, prediction}
  end

  def parse_date(string) when is_binary(string) do
    [month, day, year] = parse_tokens(string)
    {day,_} = Integer.parse(day)
    {year,_} = Integer.parse(year)
    month = month_number(month)
    {year, month, day}
  end

  defp parse_tokens(string) when is_binary(string) do
    string
    |> String.split(~r/\s/)
    |> Enum.map(&String.replace(&1, ~r/[^\w]/, ""))
  end

  @months ~w[january february march april
             may june july august september
             october november december]
             |> Enum.with_index
             |> Enum.map(fn {month, number} -> {month, number + 1} end)
  defp month_number(month) do
    {_,number} = @months |> List.keyfind(String.downcase(month), 0, {nil, nil})
    number
  end
end