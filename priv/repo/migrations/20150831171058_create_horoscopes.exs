defmodule Horoscope.Repo.Migrations.CreateHoroscopes do
  use Ecto.Migration

  def change do
    create table(:horoscopes, primary_key: false) do
      add :id, :uuid, primary_key: true, null: false

      add :sign, :string, null: false
      add :prediction, :text, null: false
      add :date, :date, null: false
      add :week, :integer, null: false
      add :year, :integer, null: false

      timestamps
    end

    create unique_index(:horoscopes, ~w[year week sign]a)
  end
end
