Feature: Checkout multiple items from the shopping cart with total price

  Scenario Outline: Successfully checkout multiple items from the cart
    Given I have added the following items to my cart and proceed to checkout:
      | itemID      | quantity     |
      | <item_id_1> | <quantity_1> |
      | <item_id_2> | <quantity_2> |
      | <item_id_3> | <quantity_3> |
    Then the checkout should be successful
    And  the total price should be "<total_price>"
    Examples:
      | item_id_1 | quantity_1 | item_id_2 | quantity_2 | item_id_3 | quantity_3 | total_price |
      | 1021      | 2          | 1101      | 3          | 1110      | 1          | 33000       |
      | 20003     | 1          | 1204      | 2          | 1145      | 3          | 29000       |

  Scenario: Attempt to checkout with an empty cart
    Given I have not added any items to my cart
    When I proceed to checkout
    Then I should receive an error message "cart doesn't contain any items"
