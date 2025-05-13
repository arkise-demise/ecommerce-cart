Feature: Shopping Cart Remove All Items Functionality

  Scenario Outline: Removing all items from the cart
    Given I have the following items in my cart:
      | ProductID | Quantity |
      | <item1_id> | <item1_quantity> |
      | <item2_id> | <item2_quantity> |
      | <item3_id> | <item3_quantity> |
      | <item4_id> | <item4_quantity> |
    When I remove all items from the cart
    Then I should receive a confirmation message "<confirmation_message>"
    And The cart should be empty

  Examples:
    | item1_id | item1_quantity | item2_id | item2_quantity | item3_id | item3_quantity | item4_id | item4_quantity | confirmation_message            |
    | 1101     | 1              | 1110     | 1              | 10010    | 1              | 10001      | 1              | All items removed from the cart |
    | 1010     | 1              | 1006     | 1              | 1007     | 1              | 1008       | 1              | All items removed from the cart |
