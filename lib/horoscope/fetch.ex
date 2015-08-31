defmodule Horoscope.Fetch do
  use Towel

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
  end

  defp request(url) do
    HTTPoison.get(url)
    |> fmap(&Map.get(&1, :body))
    |> fmap(&Floki.parse/1)
  end

  @dates %{
    "Capricorn"   => {{12, 22}, {1,  19}},
    "Aquarius"    => {{1,  20}, {2,  18}},
    "Pisces"      => {{2,  19}, {3,  20}},
    "Aries"       => {{3,  21}, {4,  19}},
    "Taurus"      => {{4,  20}, {5,  20}},
    "Gemini"      => {{5,  21}, {6,  20}},
    "Cancer"      => {{6,  21}, {7,  22}},
    "Leo"         => {{7,  23}, {8,  22}},
    "Virgo"       => {{8,  23}, {9,  22}},
    "Libra"       => {{9,  23}, {10, 22}},
    "Scorpio"     => {{10, 22}, {11, 21}},
    "Sagittarius" => {{11, 22}, {12, 21}}
  }
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

    {sign, Map.get(@dates, sign), prediction}
  end
end