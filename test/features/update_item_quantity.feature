Feature: Update item quantity in the cart

  Scenario Outline: Update the quantity of an item in my cart successfully
    Given I have added the following items to my cart
      | ItemID | Quantity |
      | 20003  | 1        |
      | 1145   | 2        |
    When I update the item with ID <ItemID> to quantity <Quantity>
    Then I should see a confirmation message "Item's quantity is updated!"
    And the cart should contain <Quantity> units of the item with ID <ItemID>

    Examples:
      | ItemID | Quantity |
      | 20003  | 3        |
      | 1145   | 5        |

  Scenario Outline: Fail to update item quantity with invalid data
    When I update the item with ID <ItemID> to quantity <Quantity>
    Then I should see an error message <ErrorMessage>

    Examples:
      | ItemID | Quantity | ErrorMessage                          |
      | 20003  | 0        | "Quantity must be greater than zero"  |
      | 99     | 1        | "Item not found in the cart"          |
