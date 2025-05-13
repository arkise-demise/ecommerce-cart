Feature: Shopping Cart Remove Item Functionality

  Scenario Outline: Removing an item from the cart
    Given I have the following items in my cart:
      | ProductID        | Quantity        |
      | <cart_productID> | <cart_quantity> |
    When I <action> the item with ProductID <remove_productID>
    Then I should see <result_message>

  Examples:
    | cart_productID | cart_quantity | action        | remove_productID | result_message              |
    | 1001           | 2             | remove        | 1001             | "Item removed from cart"    |
    | 1101           | 1             | try to remove | 1001             | "cart item not found"       |
