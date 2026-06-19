flowchart LR
  BUYER["Покупатель\nRedeem.io"]

  MARKETPLACE["g2a · eneba\n(покупка)"]

  REDEEM["Redeem.io\n— — —\n1. Выбор сервиса (Steam и др.)\n2. Выбор номинала ($10 и др.)\n3. Ввод email для доставки"]

  GENGINE["G-engine\n— — —\nbuy Steam code\nОплата с баланса\nRedeem.io в G-engine"]

  DELIVERY["Redeem.io\n— — —\nФормирование email\nс кодом для клиента\nОтправка email"]

  BUYER -->|"Фиатные деньги\n(Apple Pay · PayPal · Card)"| MARKETPLACE

  MARKETPLACE -->|"Digital Code / AAAcode\nВыдаёт селлер / из файла\n/ по API Redeem.io"| REDEEM

  REDEEM -->|"API G-engine\nbuy запрос"| GENGINE

  GENGINE -->|"API G-engine\nПолучение кода\nот G-engine"| DELIVERY

  DELIVERY -->|"Получение email\nс кодом"| BUYER
