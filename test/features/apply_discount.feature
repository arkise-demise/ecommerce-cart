Feature: Apply Discount on Cart Items

  Scenario Outline: Apply a discount to an item in the cart
    Given I have added the following items to my cart:
      | ItemID   | Quantity   |
      | <ItemID> | <Quantity> |
    When I apply a "<DiscountType>" discount of "<Discount>" to the item "<ItemID>" in the cart
    Then I should receive a success message "<SuccessMessage>"
    And the new price of the item should be "<NewPrice>"

    Examples:
      | ItemID | Quantity | DiscountType | Discount | SuccessMessage                                  | NewPrice |
      | 10001  | 1        | percentage   | 10       | 10.00 percentage discount applied to item 10001 | 4950.00   |
      | 1101   | 2        | flat         | 50       | 50.00 flat discount applied to item 1101        | 5450.00   |
      | 1110   | 1        | percentage   | 5        | 5.00 percentage discount applied to item 1110   | 5225.00   |
