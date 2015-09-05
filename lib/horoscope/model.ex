defmodule Horoscope.Model do
  use Ecto.Model

  @required_fields ~w[sign prediction date week year]a
  @optional_fields ~w[]a

  @primary_key {:id, :binary_id, autogenerate: true}
  schema "horoscopes" do
    field :sign, :string
    field :prediction, :string
    field :date, Ecto.Date
    field :week, :integer
    field :year, :integer

    timestamps
  end

  def changeset(model, params \\ :empty) do
    model
    |> cast(params, @required_fields, @optional_fields)
    |> update_change(:sign, &String.downcase/1)
    |> unique_constraint(:sign, name: :horoscopes_year_week_sign_index)
  end
end