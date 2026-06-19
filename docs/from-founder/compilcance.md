Compliance-проверка Eneba
1. Статус площадки для нас
Eneba подходит как приоритетная площадка для MVP, потому что это gaming marketplace с поддержкой цифровых товаров, gift cards и merchant/API-инфраструктуры. У Eneba есть отдельная страница для продавцов, где они заявляют merchant-панель, отчёты, расчёты и API-интеграцию для автоматизации продаж. 
Предварительный вывод: Eneba подходит как целевой marketplace, но вход будет зависеть от того, сможем ли мы доказать легальный источник gift cards / цифровых кодов и показать безопасный процесс выдачи.

---
2. Требования к продавцу
Eneba не выглядит как площадка с открытой самостоятельной регистрацией любого продавца. На странице Sell on Eneba указано, что нужно связаться с ними через форму, после чего заявка проходит оценку, и представитель Eneba сообщает детали подключения. API также доступен только мерчантам, продающим цифровые продукты. 
Что это значит для MVP: нам нужно готовить не просто лендинг, а пакет для merchant approval:
-  описание компании / юрлица; 
-  описание продукта; 
-  источник gift cards / цифровых кодов; 
-  процесс выдачи товара пользователю; 
-  refund / delivery / support policy; 
-  антифрод; 
-  proof of delivery; 
-  marketplace-ready страницу покупки/активации. 

---
3. Source of stock — главный блокер
Самое важное требование Eneba: продавцы должны получать товары от официальных дистрибьюторов. Eneba прямо пишет, что товары, купленные у других продавцов Eneba, не предназначены для перепродажи, и что продавать могут только мерчанты, sourcing products from official distributors. При необходимости Eneba может запросить proof of purchase, чтобы проверить, что stock получен от authorized wholesale distributor. 
Это ключевой риск для нашего MVP.
Что нужно подготовить:
-  договор / подтверждение от Jinjin / Jeenjean или другого поставщика; 
-  invoices / proof of purchase; 
-  описание supplier chain; 
-  список SKU и номиналов; 
-  подтверждение, что коды не куплены на другом marketplace; 
-  процесс проверки проблемных кодов через supplier support. 
Вывод: без подтверждённого source of stock шанс пройти Eneba низкий.

---
4. Модель external redeem
Наш MVP предполагает, что пользователь покупает ваучер на marketplace, а затем активирует его на нашем сайте и получает цифровой товар. Нужно отдельно подтвердить у Eneba, допустима ли такая модель.
С одной стороны, Eneba поддерживает продажу digital offers и merchant API. Также в API есть логика автоматизации stock, где marketplace может обращаться к системе продавца для reservation/provision: сначала резервирует ключ, потом запрашивает выдачу купленного ключа. 
С другой стороны, для Eneba может быть чувствительно, если покупатель получает не конечный товар сразу, а универсальный ваучер, который нужно активировать вне площадки. Это может восприниматься как дополнительный redemption layer.
Что нужно проверить напрямую у Eneba:
-  можно ли продавать universal voucher, который редимится на внешнем сайте; 
-  можно ли вместо конечного Steam/Xbox key выдавать voucher code; 
-  разрешён ли внешний /redeem flow; 
-  нужно ли, чтобы конечный цифровой товар выдавался внутри Eneba delivery flow; 
-  можно ли использовать Declared Stock / Provision API для выдачи нашего ваучера; 
-  допускаются ли инструкции “перейдите на наш сайт и активируйте код”. 
Предварительный verdict: модель возможна, но требует подтверждения Eneba. Самый безопасный вариант для acceptance — сделать так, чтобы Eneba видела это как продажу конкретного digital voucher/gift product, а не как финансовый инструмент.

---
5. API и автоматизация
Eneba предоставляет GraphQL API для мерчантов цифровых товаров. API защищён OAuth 2.0 и доступен только approved merchants. 
Для нас особенно важен механизм Declared Stock: он позволяет не загружать ключи заранее в Eneba, а объявить наличие товара. При покупке Eneba вызывает систему продавца дважды: сначала reservation request, затем provision request для выдачи ключей. 
Это хорошо ложится на наш MVP:
-  партнёр / мы объявляем stock; 
-  Eneba резервирует заказ; 
-  наша система выдаёт ваучер или конечный код; 
-  можно контролировать availability через API; 
-  меньше ручной загрузки кодов. 
Но это также означает, что для MVP нужно обеспечить:
-  стабильный endpoint для reserve/provision; 
-  быстрый ответ; 
-  обработку ошибок; 
-  корректные статусы; 
-  защиту от двойной выдачи; 
-  мониторинг failed provision. 

---
6. Product listing и новые товары
Eneba позволяет создавать digital offers через Vendor Panel. Если товара нет в каталоге, можно запросить добавление нового product name через merchant account, после чего команда Eneba рассмотрит заявку. 
Для нас это значит:
-  если мы продаём “Brand Name Voucher”, возможно, понадобится отдельная категория/товар; 
-  для universal voucher может потребоваться approval нового продукта; 
-  нужно заранее подготовить название, описание, инструкции, ограничения и terms; 
-  важно не использовать слова, которые делают продукт похожим на wallet/balance/payout. 
Рекомендованное позиционирование для Eneba:
Digital gaming gift voucher redeemable for selected gaming gift cards and digital goods.
Не использовать:
-  wallet; 
-  balance; 
-  payout; 
-  cashout; 
-  crypto withdrawal; 
-  bank card; 
-  prepaid balance. 

---
7. Gift card terms и риск resale / transfer for value
У Eneba есть отдельные условия для gift cards, где указано, что gift card balance не может перепродаваться, передаваться за ценность или использоваться для unauthorized commercial purposes. 
Это относится к Eneba gift card balance, но для нас важно как сигнал: Eneba чувствительна к сценариям, где gift card превращается в transferable value.
Чтобы снизить риск:
-  ваучер должен использоваться полностью; 
-  не должно быть баланса; 
-  не должно быть остатка; 
-  не должно быть перевода между пользователями; 
-  не должно быть cash-like redemption; 
-  пользователь должен получать конкретный цифровой товар. 
Это полностью совпадает с текущим ТЗ и усиливает no-wallet/no-balance позиционирование.

---
8. Disputes, refunds и proof of delivery
Eneba как marketplace фокусируется на конечном пользователе, поэтому disputes по неработающим / уже использованным ключам будут критичным риском. Публичные user complaints по marketplace keys часто связаны именно с already redeemed / invalid codes. Это не официальный источник требований, но показывает типовую боль категории. 
Для MVP обязательно нужно иметь:
-  delivery timestamp; 
-  voucher hash; 
-  SKU; 
-  supplier transaction ID; 
-  delivery page event; 
-  email delivery log; 
-  статус Delivered; 
-  proof of delivery для партнёра; 
-  процедуру проверки “код уже использован”. 
Иначе продавец будет слаб в спорах с покупателем и площадкой.

---
9. Delivery SLA
У Eneba есть API-модель, в которой ключ должен быть выдан при provision request. Это значит, что задержки и pending-сценарии должны быть минимальны. 
Для MVP лучше целиться в:
-  обычная выдача: до 30 секунд; 
-  максимум: до 2 минут; 
-  если supplier не отвечает — товар лучше скрывать заранее; 
-  долгие продукты не включать в MVP; 
-  все ошибки provision должны попадать в alert. 
Если мы оставляем сценарий “5–15 минут”, он может быть приемлем для нашего сайта, но для Eneba acceptance может выглядеть хуже, если покупатель ожидает мгновенную выдачу.

---
10. Что нужно подготовить для общения с Eneba
Минимальный пакет:
 Описание продукта
 Что такое voucher, как работает, что получает пользователь. 
Source of stock
 Supplier agreement, invoices, API docs, proof от официального/оптового поставщика. 
UX flow
 Скриншоты /how-to-buy, /redeem, /delivery. 
Terms / Policies
 Terms, Privacy, Refunds, Redemption & Delivery, Support. 
Anti-fraud
 Rate limits, captcha, hash storage, brute-force protection. 
Proof of Delivery
 Формат отчёта для споров. 
API readiness
 Поддержка reserve/provision или альтернативный delivery flow. 
Support SLA
 Как обрабатываются failed/invalid/already redeemed cases. 

---
11. Blockers для MVP
Критические blockers
Нет подтверждённого source of stock от официального supplier. 
Eneba не разрешает external redeem для universal voucher. 
Продукт выглядит как wallet/balance/cash-like instrument. 
Нет стабильной мгновенной выдачи. 
Нет proof of delivery. 
Средние риски
Нужно добавлять новый product/category. 
Ограничения по регионам SKU. 
Высокие требования к seller history. 
Hold/verification на merchant onboarding. 
Высокая чувствительность к complaints/disputes. 

---
12. Verdict по Eneba
Предварительный статус: PASS WITH LIMITATIONS.
Eneba подходит для MVP, но только если выполнить 4 условия:
1.  Подтвердить официальный source of stock. 
2.  Согласовать external redeem / universal voucher model. 
3.  Убрать любые признаки wallet/balance/payout. 
4.  Обеспечить быстрый fulfillment и proof of delivery. 
Если Eneba не принимает external redeem, альтернативный вариант — выдавать через Eneba конкретный digital code/SKU с автоматической provision-логикой через API.

Compliance-проверка Kinguin
1. Статус площадки для нас
Kinguin подходит для MVP как крупная gaming-площадка с фокусом на цифровые товары: game keys, activation codes, gift cards и другие digital goods. В их legal notice прямо указано, что товары на Kinguin продаются в виде уникальных ключей, которые покупатель активирует на сторонних платформах, например Steam. Это хорошо совпадает с нашей логикой ваучер - цифровой товар, если продукт не выглядит как wallet/payout. 
Kinguin также имеет развитую API-инфраструктуру для продавцов: API позволяет автоматизировать процесс от создания offer до отправки ключа покупателю. 
Kinguin потенциально подходит для MVP, особенно как площадка, где важна автоматизация, цифровые ключи и merchant API.

---
2. Требования к продавцу / merchant onboarding
Kinguin работает с продавцами через модель verified merchant. Сейчас нельзя свободно продавать ключи как Community Seller: для продажи нужно подать заявку через страницу Sell on Kinguin и пройти проверку мерчанта. Это подтверждается support-разделом Kinguin: они прямо указывают, что для продажи нужно стать verified merchant, а формат Community Seller сейчас недоступен. 
Для MVP это значит, что нам нужно готовить полноценную merchant заявку. Также важно учитывать, что доступ к API и merchant-функциям появляется после создания и подтверждения merchant account. В API guide Kinguin описывает процесс создания merchant account через Sell on Kinguin и дальнейший переход в merchant account после проверки и approval. 
Минимально нужно подготовить:
-  данные компании; 
-  контактное лицо; 
-  описание продукта и категорий товаров; 
-  источник цифровых кодов / gift cards; 
-  подтверждение легальности stock; 
-  примерные объёмы продаж; 
-  ссылку на сайт / landing page; 
-  описание процесса выдачи товара; 
-  support и refund policy. 
Дополнительно в условиях Kinguin указано, что предприниматели/компании могут быть обязаны предоставить адрес, телефон, банковские данные, адрес бизнеса, налоговый номер и регистрационный номер компании. Kinguin также оставляет за собой право проверять данные продавца и запрашивать документы, подтверждающие адрес компании, регистрацию, право представлять компанию и VAT/tax ID. 
Вывод: для Kinguin важно заранее подготовить merchant package: юридические данные, описание продукта, подтверждение источника кодов, сайт, support/refund policy и схему выдачи товара. Без этого подключение может затянуться на этапе проверки или approval.

---
3. Source of stock / происхождение ключей
Для Kinguin это один из ключевых рисков. Площадка исторически работает с цифровыми ключами и game activation codes, поэтому вопросы происхождения stock, валидности ключей и disputes будут критичны.
В их terms указано, что Kinguin предоставляет площадку для продажи games and game activation codes пользователям со стороны sellers. 
Для нашего MVP нужно подготовить:
-  supplier agreement / договор с Jinjin или другим поставщиком; 
-  invoices / подтверждение закупки; 
-  описание supplier chain; 
-  список SKU; 
-  номиналы и регионы; 
-  процесс проверки invalid / already used codes; 
-  proof of delivery; 
-  support flow для спорных ситуаций. 
Критичный вопрос к Kinguin
Можно ли продавать не конечный Steam/Xbox key сразу, а универсальный ваучер, который активируется на нашем сайте и затем превращается в gift card / digital good?
Это нужно подтвердить отдельно, потому что Kinguin привычнее к модели продажи конкретного digital key, а не внешнего redeem-слоя.

---
4. External redeem model
Kinguin уже допускает цифровые ключи, которые активируются на сторонних платформах. Их legal notice описывает товары как уникальные ключи для redeem на third-party platforms. 
Но наш ваучер промежуточный universal voucher. Поэтому Kinguin может рассматривать его как:
-  digital activation code — допустимо; 
-  gift voucher — допустимо при правильном описании; 
-  stored value / wallet-like product — рискованно. 
Как лучше позиционировать
-  digital gaming gift voucher; 
-  redeemable for selected gaming gift cards; 
-  activation code; 
-  no balance; 
-  full redeem only; 
-  no cash-out. 

---
5. API и автоматизация
Kinguin сильная площадка именно с точки зрения API. Их developer portal описывает API для управления merchant-процессами, stock, offers, wholesale, reservations и alerts. 
Также Kinguin API позволяет продавцам обрабатывать весь процесс сделки через API: от создания offer до отправки ключа покупателю. 
Для доступа к API нужно включить 2FA, создать API client и получить secret key. 
Для MVP это хорошо: наш продукт можно строить как API-first систему.
Что нужно подготовить технически:
-  генерация voucher codes; 
-  загрузка/синхронизация stock; 
-  endpoint или процесс для выдачи ключа; 
-  логика “код выдан только один раз”; 
-  хранение статусов; 
-  обработка отмены/ошибок; 
-  мониторинг failed issuance; 
-  support tools. 

---
6. Product listing и категории
Kinguin специализируется на игровых ключах, software, gift cards и digital goods. Для нас важно понять, можно ли создать отдельный продукт формата: [Brand Name] Gaming Gift Voucher $10 / $25 / $50 / $100
или потребуется размещать конкретные SKU: Steam Gift Card $10, Xbox Gift Card $25, PSN Gift Card $50
Риск в том, что universal voucher может не вписаться в стандартный listing flow. Если так, лучше начать с конкретных SKU или брендированных ваучеров с понятным назначением.
Что нужно уточнить у Kinguin
-  можно ли продавать universal gaming voucher; 
-  можно ли использовать external redeem page; 
-  можно ли в листинге давать инструкцию активируйте на нашем сайте; 
-  можно ли продавать ваучер, который потом пользователь обменивает на несколько gift cards; 
-  можно ли создавать собственный брендовый voucher; 
-  какие категории доступны для такого товара. 

---
7. Delivery expectations
Kinguin покупатели ожидают быструю выдачу цифрового ключа. Поскольку API поддерживает процесс создание offer - отправка ключа покупателю, MVP должен обеспечивать максимально короткий delivery time. 
Рекомендуемый SLA для MVP:
-  нормальная выдача: до 30 секунд; 
-  максимум: до 2 минут; 
-  долгие товары не включать в стартовый каталог; 
-  если supplier API нестабилен — SKU скрывать или отключать; 
-  обязательно иметь alerts при failed issuance. 

---
8. Disputes, refunds и support
Для Kinguin критично минимизировать кейсы:
-  код не работает; 
-  код уже использован; 
-  товар не получен; 
-  непонятно, где активировать. 
Kinguin как marketplace с цифровыми ключами чувствителен к качеству продавца, жалобам и репутации. Поэтому нужно подготовить:
-  Proof of Delivery; 
-  delivery timestamp; 
-  voucher hash; 
-  email delivery log; 
-  supplier transaction ID; 
-  статус Delivered; 
-  инструкции по использованию; 
-  support SLA; 
-  процесс проверки already used codes. 
Для MVP обязательно
Партнёр должен иметь возможность быстро показать, что ваучер был активирован, цифровой товар был выдан, операция завершена.

---
9. Compliance wording
Для Kinguin критично не выглядеть как платёжный или cash-like продукт.
Пример безопасного описания:
Digital gaming gift voucher that can be activated on [domain] and exchanged for selected gaming gift cards and digital products. No balance, no cash-out, full redeem only.

---
10. Финансовые и payout-риски
Kinguin как marketplace будет платить продавцу по своей стандартной схеме. Для нас важно не смешивать это с нашей внутренней логикой.
Риск появляется, если:
-  мы начинаем предлагать конечному пользователю payout; 
-  появляется wallet; 
-  пользователь может хранить value; 
-  есть crypto withdrawal; 
-  есть остаток после redeem. 
Поэтому для Kinguin-версии MVP нужно строго сохранить:
full redeem only
no balance
no payout
digital goods only

---
11. Что нужно подготовить для общения с Kinguin
Минимальный пакет:
 Описание продукта
 Что такое voucher, как работает, что получает пользователь. 
 Список SKU
 Номиналы, регионы, доступные товары. 
 Source of stock
 Supplier agreement, invoices, API docs, proof of legitimate sourcing. 
 UX flow
 Скриншоты landing / redeem / delivery page. 
 Delivery policy
 SLA: обычно до 30 секунд, максимум до 2 минут. 
 Refund / dispute policy
 Что делаем с invalid / already used / not received. 
 Proof of Delivery
 Формат отчёта для споров. 
 API readiness
 Как выдаётся и отслеживается ключ. 
 Support flow
 Как пользователь и marketplace получают помощь. 

---
12. Blockers для MVP
Критические blockers
 Kinguin не разрешит продавать universal voucher с external redeem. 
 Нет подтверждённого source of stock. 
 Продукт выглядит как wallet / stored value. 
 Нет быстрой выдачи товара. 
 Нет доказательства доставки. 
 Нет verified merchant account. 
Средние риски
 Нужно создавать новую категорию/продукт. 
 Высокая конкуренция по цене. 
 Жалобы покупателей на непонятный redeem flow. 
 Ограничения по регионам gift cards. 
 Требования к seller history. 

---
13. Verdict по Kinguin
Предварительный статус: PASS WITH LIMITATIONS.
Kinguin подходит для MVP, если:
-  Получить verified merchant account. 
-  Подтвердить легальный source of stock. 
-  Согласовать модель external redeem. 
-  Обеспечить быструю выдачу digital goods. 
-  Убрать любые wallet/payout/cash-like признаки. 
-  Подготовить proof of delivery и support flow. 
Если Kinguin не разрешит universal voucher, fallback-сценарий - продавать конкретные branded digital gift products с выдачей через нашу систему.
