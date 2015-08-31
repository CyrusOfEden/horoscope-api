defmodule Horoscope.Scraper do
  use Towel
  alias Horoscope.Model
  alias Horoscope.Repo

  def seed do
    fetch
    |> Stream.with_index
    |> Stream.flat_map(fn {horoscopes, offset} ->
      iso_week = week_number(offset)
      Enum.map(horoscopes, &changeset(&1, iso_week))
    end)
    |> Stream.map(&Repo.insert/1)
  end

  @base "http://www.theonion.com/features/horoscope"
  def fetch do
    request(@base)
    |> fmap(&Floki.find(&1, ".reading-list-item"))
    |> fmap(&Floki.attribute(&1, "data-absolute-url"))
    |> fmap(&Stream.map(&1, fn item -> get(item) end))
    |> Result.unwrap
  end

  defp get(page) do
    request(page)
    |> fmap(&Floki.find(&1, ".astro"))
    |> fmap(&Floki.find(&1, ".large-thing"))
    |> fmap(&Enum.map(&1, fn elem -> elem |> Floki.text |> normalize end))
    |> Result.unwrap
  end

  defp changeset({sign, prediction}, {year, week}) do
    Model.changeset(%Model{}, %{
      year: year,
      week: week,
      sign: sign,
      prediction: prediction
    })
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

  def week_number(offset) do
    {year, week} = :calendar.iso_week_number()
    if week > offset do
      {year, week - offset}
    else
      {nil, nil}
    end
  end
end