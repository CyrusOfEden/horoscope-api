defmodule Horoscope.Model do
  use Ecto.Model

  @required_fields ~w[week year sign prediction]a
  @optional_fields ~w[]a

  @primary_key {:id, :binary_id, autogenerate: true}
  schema "horoscopes" do
    field :year, :integer
    field :week, :integer
    field :sign, :string
    field :prediction, :string

    timestamps
  end

  def changeset(model, params \\ :empty) do
    model |> cast(params, @required_fields, @optional_fields)
  end
end