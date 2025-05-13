Feature: Shopping Add Cart Functionality

Scenario Outline: Adding items to the cart
    Given I have the following products in stock:
        | ProductID | Quantity |
        | <ProductID> | <Quantity> |
    When I add these items to my cart
    Then I should receive "items successfully added to the cart"

    Examples:
    | ProductID | Quantity |
    | 1101      | 1        |
    | 1110      | 1        |

Scenario Outline: Adding more than 10 unique items
    Given I already have the following unique items in my cart:
        | ProductID     | Quantity      |
        | <ProductID1>  | <Quantity1>   |
        | <ProductID2>  | <Quantity2>   |
        | <ProductID3>  | <Quantity3>   |
        | <ProductID4>  | <Quantity4>   |
        | <ProductID5>  | <Quantity5>   |
        | <ProductID6>  | <Quantity6>   |
        | <ProductID7>  | <Quantity7>   |
        | <ProductID8>  | <Quantity8>   |
        | <ProductID9>  | <Quantity9>   |
        | <ProductID10> | <Quantity10>  |
    When I try to add one more unique item to the cart:
        | ProductID     | Quantity     |
        | <ProductID11> | <Quantity11> |
    Then I should see an error "can't add more than 10 unique items to the cart"

    Examples:
    | ProductID1 | Quantity1 | ProductID2 | Quantity2 | ProductID3  | Quantity3 | ProductID4  | Quantity4 | ProductID5 | Quantity5 | ProductID6 | Quantity6 | ProductID7 | Quantity7 | ProductID8 | Quantity8 | ProductID9 | Quantity9 | ProductID10 | Quantity10 | ProductID11 | Quantity11 |
    | 1101       | 1         | 1110       | 1         | 10010       | 1         | 10001       | 1         | 1011       | 1         | 1006       | 1         | 1007       | 1         | 1008       | 1         | 1009       | 1         | 1010        | 1          | 1100          | 1        |
