Feature: View items in the cart and their total price

  Scenario Outline: View the items added to the cart and get the total price
    Given I add the following items to my cart and view them:
      | item_id     | quantity     | 
      | <item_id_1> | <quantity_1> | 
      | <item_id_2> | <quantity_2> | 
      | <item_id_3> | <quantity_3> | 
    Then I should see the following items:
      | item_id     | item_name     | quantity     | price     |
      | <item_id_1> | <item_name_1> | <quantity_1> | <price_1> |
      | <item_id_2> | <item_name_2> | <quantity_2> | <price_2> |
      | <item_id_3> | <item_name_3> | <quantity_3> | <price_3> |
    And the total price of the cart should be "<total_price>"

  Examples:
    | item_id_1 | item_name_1   | quantity_1 | price_1 | item_id_2 | item_name_2   | quantity_2 | price_2 | item_id_3 | item_name_3   | quantity_3 | price_3 | total_price |
    | 1021      | noxiabag      | 2          | 5500    | 1101      | jackasogets   | 3          | 5500    | 1110      | bostenoeven   | 1          | 5500    | 33000       |
    | 20003     | cotex-shirt   | 1          | 12000   | 1204      | Nike-Shoes    | 2          | 4000    | 1145      | Nike-Tshirt   | 3          | 3000    | 29000       |
