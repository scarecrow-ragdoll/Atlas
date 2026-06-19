flowchart LR
  AAA["Redeem.io"]

  SELLER["Продавец Redeem.io\nDigital Marketplace"]

  BUYER["Покупатель\ng2a · eneba · etc"]

  SUPPLIER["Поставщик\nG-engine · Digital Code"]

  SELLER -->|"1. Пополнение баланса\n(Crypto)"| AAA
  AAA -->|"2. Получает Redeem.io's"| SELLER
  AAA -->|"0. Пополнение баланса\nпоставщика (Crypto)"| SUPPLIER
  BUYER -->|"3. Оплата за Redeem.io"| SELLER
  SELLER -->|"4. Выдача Redeem.io"| BUYER
  BUYER -->|"5. Redeem AAACode\nВыбор Digital Code"| SUPPLIER
  SUPPLIER -->|"6. Выдача кода\n(Steam и др.)"| BUYER
