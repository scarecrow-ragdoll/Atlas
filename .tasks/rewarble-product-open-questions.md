# Дополнительные вопросы к продуктовой спецификации Rewarble/Redeem.io

Статус: актуализировано 2026-05-27 по `docs/from-founder/*`.

Назначение: этот документ фиксирует только те уточнения, которые еще нужны команде реализации, чтобы спроектировать и собрать MVP без догадок. Вопросы, на которые founder-доки уже дали рабочий ответ, вынесены в раздел "Закрытые вопросы" и больше не должны считаться открытыми, пока не изменится источник.

Адресат открытых вопросов: заказчик, product owner, compliance/legal, finance, operations или технический владелец интеграции. Если ответ зависит от провайдера, маркетплейса, юристов, финансов или операций, ожидаемое действие - предоставить утвержденный документ, контракт, API-документацию, доступ, контакт ответственного или финальное решение.

## Источники синхронизации

- `docs/from-founder/prd.md`
- `docs/from-founder/tz.md`
- `docs/from-founder/compilcance.md`
- `docs/from-founder/analytics.md`
- `docs/from-founder/core-flow.md`
- `docs/from-founder/role-details.md`
- `docs/from-founder/reedem.io.md`
- `docs/from-founder/seller-flow.md`
- `docs/from-founder/supplier-flow.md`
- `docs/from-founder/start-of-sku.md`
- `docs/from-founder/dealers-schema.md`

## Закрытые вопросы

Список ниже оставлен для навигации по старым номерам. Это уже не открытые вопросы; если в будущем founder-доки поменяются, такие пункты можно вернуть в открытые.

- Старый 1: базовая модель MVP - дистрибьюторская маржа. Партнер покупает ваучеры у платформы с дисконтом, продает на маркетплейсе, платформа зарабатывает между ценой партнеру и закупкой у поставщика. Исключения и точная экономика вынесены в открытые вопросы.

- Старый 6: "депозит" - B2B предоплата/доступный лимит партнера, не B2C-кошелек. Выпуск ваучеров только по предоплате; доступный лимит считается как `confirmed deposits - spent`.

- Старый 7: запрещены wallet/balance/payout/withdrawal/cash-like формулировки и визуальные паттерны. Допустимые рабочие термины: "доступный лимит/депозит", "B2B предоплата".

- Старые 12, 13: стартовый каталог-черновик задан: Steam, PlayStation, Xbox, Nintendo, Riot Games, Roblox; Fortnite/Epic Games не раскрыт. Базовая матрица `SKU -> номинал -> регион` задана для стартовых SKU. Country-level restrictions и трактовка `Global / US` остаются открытыми.

- Старый 15: ваучер задуман как универсальный gift code. Пользователь выбирает digital-товар из каталога, но нет баланса, payout или остатка.

- Старый 16: корзина разрешена только если сумма выбранных товаров ровно равна номиналу ваучера; больше/меньше нельзя, остаток всегда 0.

- Старый 21: источники конечных кодов - manual stock и API supplier. MVP Pilot/P1 стартует с ручным stock; Jinjin API live и автоматический stock sync вынесены в P2/Scale.

- Старый 29: срок действия ваучера - 12 месяцев; после истечения статус `EXPIRED`, редим невозможен.

- Старый 30: ваучер считается потраченным/consumed только при переходе в `DELIVERED`; в `PENDING_FULFILLMENT` он зарезервирован под заказ и не может быть использован повторно.

- Старый 31: если сбой случился после подтверждения, заказ остается `PENDING_FULFILLMENT`; пользователю предлагаются альтернативы того же номинала, support может заменить SKU внутри заказа, сохраняя остаток 0.

- Старый 34: основной SLA-контур задан: P95 fulfillment <= 30 секунд, максимальный timeout 2 минуты, delayed email после 15 минут, manual escalation после 6 часов, retry schedule до 6 часов. Конфликт с текстом "обычно 5-15 минут" вынесен в открытые вопросы.

- Старые 43, 110: B2C support разделен. Вопросы оплаты или неполучения кода от продавца идут к продавцу/маркетплейсу; если код не активируется на нашем сайте - к support Redeem.io.

- Старый 44: онбординг партнера - регистрация и KYB-заявка со статусами; админ затем approve/suspend и управляет лимитами.

- Старый 53: минимальные статусы batch - `created`, `delivered_to_partner`, `cancelled`.

- Старые 55, 56: базовый состав PoD - `voucher_hash`, `batch_id`, timestamp редима, SKU/номинал, статус `Delivered`, internal/supplier tx id; без PII конечного пользователя и без раскрытия назначения. Процесс спора и юридическая достаточность остаются открытыми.

- Старый 58: Partner API - создать ваучер(ы)/batch; API keys + rotation, idempotency keys, rate limits, signed status webhooks. По PRD это P2/Tier 1+, конфликт с TZ `MUST` вынесен в открытые вопросы.

- Старые 60, 61: минимальные внутренние роли - Support видит статусы/masked email/код скрыт; Ops может resend/retry, batch counters и reveal once с логом; Admin - все плюс лимиты и экспорт.

- Старые 65, 66, 67: Admin/Ops order tooling - поиск по `order_id / email / voucher_hash`, карточка заказа с timeline/fulfillment attempts/masked code; действия resend email, retry delivery/fulfillment, switch supplier, сменить SKU, mark resolved.

- Старые 74, 113: замена SKU разрешена только на эквивалент по номиналу; итоговая сумма = номинал ваучера, остаток 0, баланс не создается.

- Старый 78: минимальный legal set перечислен: Terms of Use, Redemption & Delivery Policy, Refunds & Cancellations, Privacy Policy, Cookie Policy, Sanctions & Restricted Use, Partner Terms. Также есть Partner Compliance, Support & Complaints и PoD Policy.

- Старые 100, 115: публичные B2C-экраны MVP перечислены: `/`, `/how-to-buy`, `/redeem`, `/delivery/:token`, `/redeem/error`, `/faq`; `/how-to-buy` показывает G2A, Kinguin, Eneba.

- Старый 104: базовые сущности продукта перечислены: Partner, VoucherBatch, Voucher, RedemptionOrder, FulfillmentItem, DeliveryToken. Связи и техническая нормализация остаются открытыми.

- Старые 107, 108, 116: граница MVP/P1 в целом задана: B1-B7, P1-P6, A1-A7, F1-F3, S1-S3, manual stock, один партнер. P2: Partner API, Jinjin API live, Eneba + Kinguin live, расширенный антифрод, automatic stock sync.

- Старый 118: финальность - ваучер consumed только при `DELIVERED`; после выдачи конечного кода операция завершена, возвраты/отмена невозможны, кроме явной ошибки системы до выдачи.

- Старый 131: базовый KYB-пакет - регистрация, директора, UBO, sanctions/PEP/adverse media, source of funds, подтверждение витрин/аккаунтов маркетплейса. Порог UBO и периодичность проверок остаются открытыми.

- Старый 132: базовая KYT-логика - high-risk или санкционные связи ведут к отказу/hold и ручной проверке. Post-confirm state machine остается открытым.

## Актуальные открытые вопросы

### A. Бизнес-модель, деньги и коммерческие правила

1. Какие исключения из дистрибьюторской маржи допустимы по партнеру, SKU, маркетплейсу или промо-условиям?
   Почему важно: pricing, invoice lines, лимиты и отчеты должны учитывать исключения, но базовая модель уже закрыта.

2. Нужно зафиксировать формулу цены для партнера по каждому ваучеру/SKU: поля расчета, уровень настройки, момент фиксации, округление, скидка и правила пересмотра.
   Почему важно: выпуск batch и списание доступного лимита требуют детерминированной цены до экспорта кодов.

3. 5-8% - это revenue, gross margin или net margin после COGS, payment costs, fraud reserve и операционных расходов?
   Почему важно: founder analytics отдельно просит разделить GMV, revenue, COGS, payment costs, fraud reserve, gross margin и net margin.

4. Есть ли в B2C redeem пользовательская комиссия, FX/service fee или экономика только B2B distribution margin?
   Почему важно: B2C UI, legal copy, receipts и dispute handling отличаются, если пользователь видит комиссию.

5. Какие валюты поддерживаются для номинала ваучера, SKU, платежей партнера, стоимости поставщика и отчетности?
   Почему важно: мультивалютность влияет на каталог, CSV, FX, PoD, reconciliation и accounting.

6. Должна ли валюта ваучера всегда совпадать с валютой SKU, или разрешен cross-currency redeem через FX?
   Почему важно: корзина, exact amount validation и маржа ломаются без правила валютного соответствия.

7. Нужно утвердить FX-policy: источник курса, момент фиксации, округление, пересчет для лимитов, маржи, PoD и споров.
   Почему важно: суммы должны воспроизводиться в отчетах и спорных случаях.

8. Какие способы оплаты партнера входят в MVP: wire, crypto, оба варианта или ручное offline-confirmation?
   Почему важно: PRD говорит `wire and/or crypto`; реализация платежей, KYT и сверка зависят от финального выбора.

9. Для crypto-платежей нужно зафиксировать сети, активы, whitelist-правила, KYT-провайдера, число подтверждений, холды и post-confirm high-risk state machine.
   Почему важно: выпуск batch после on-chain подтверждения может конфликтовать с поздним KYT-risk сигналом.

10. Как обрабатываются переплаты, недоплаты, ошибочные переводы, failed deposits, chargebacks, выход партнера и возвраты на источник?
    Почему важно: без правила невозможно согласовать отсутствие payout с возвратами и бухгалтерским учетом.

11. Где юридически и финансово хранятся предоплаты партнеров: сегрегированный счет, операционные средства, кредит поставщика, дебиторская задолженность или другая модель?
    Почему важно: это влияет на custody, лимиты, partner terms, accounting и regulatory classification.

12. Может ли партнер отменить выпущенный batch или вернуть средства за неиспользованные ваучеры?
    Почему важно: жизненный цикл batch, voucher status и поддержка партнеров должны быть конечными.

13. Какая модель инвойсов, налогов, receipts, credit notes и выгружаемых выписок нужна для депозитов и покупки batch?
    Почему важно: финансовые экраны, exports, VAT/tax fields и reconciliation не определены.

14. Какой договорный SLA и ценностное обещание партнеру реально утверждены?
    Почему важно: внутренние алерты, публичные тексты, onboarding и договоры должны совпадать.

### B. Каталог, SKU, остатки и поставщики

15. Какие конкретные SKU входят в P1 из стартового каталога, какие остаются post-MVP, и что делать с Fortnite/Epic Games?
    Почему важно: `start-of-sku.md` дает список, но не финальный P1 scope.

16. Как моделировать `Global / US`: один SKU с несколькими allowed regions или отдельные SKU/variants?
    Почему важно: региональные ограничения влияют на каталог, UX и споры по неверному региону.

17. Нужна полная матрица `SKU -> country/region/restrictions/allowed partner/marketplace`.
    Почему важно: текущий список содержит базовые регионы, но не country-level restrictions.

18. По какому правилу система выбирает каталог после ввода ваучера: partner, marketplace, region, currency, denomination, stock, risk tier или комбинация?
    Почему важно: фильтрация каталога и API-контракт `/redeem` зависят от этого.

19. Что происходит, если ни одна точная комбинация SKU не равна номиналу ваучера?
    Почему важно: full redeem only требует понятный отказ, alternative flow или support route.

20. Можно ли в одной корзине смешивать разные бренды, категории, регионы или валюты?
    Почему важно: PoD, supplier fulfillment и marketplace expectations зависят от состава корзины.

21. При multi-item redeem выдача должна быть all-or-nothing или допускается частичная доставка?
    Почему важно: при частичной доставке нужны item statuses, partial failure и отдельный PoD по позиции.

22. Нужно зафиксировать точный TTL soft reservation: 1, 2 или 3 минуты, как показывается таймер и что происходит после истечения.
    Почему важно: founder docs дают диапазон 1-3 минуты, а тестам и UX нужно одно значение.

23. Если soft reservation истекла до подтверждения, что видит пользователь и как освобождается stock?
    Почему важно: гонки остатков и повторный выбор нельзя безопасно реализовать.

24. Для каждого P1 SKU нужно указать источник: `manual_stock`, `api_supplier`, оба источника, и где обязателен ручной fallback для top SKU.
    Почему важно: P1 стартует с manual stock, но модель источников нужна на уровне SKU.

25. Какой приоритет источников по SKU: manual stock first, API supplier first, цена/надежность/скорость, partner-specific routing?
    Почему важно: fallback и счетчики остатков нельзя вывести из общей формулировки.

26. Предоставьте API-документацию Jinjin/Jeenjean/G-engine: auth, endpoints, request/response, sandbox, rate limits, SLA, error codes, idempotency и тестовые доступы.
    Почему важно: supplier integration остается блокером P1/P2.

27. Нужно определить taxonomy ошибок поставщика: retryable, fallback, final failed, manual escalation, invalid/already used supplier code.
    Почему важно: state machine fulfillment и alerting зависят от классификации ошибок.

28. Нужно определить supplier reconciliation artifacts: supplier order id, transaction id, cost, settlement status, refunds, stale/duplicate codes и dispute evidence.
    Почему важно: COGS, PoD и разбор invalid codes требуют внешних ID.

29. Нужен supplier due diligence/KYS-процесс: authorized distributor proof, invoices, sanctions checks, ongoing monitoring.
    Почему важно: marketplace approval требует легального source of stock.

30. Как хранятся и защищаются исходные конечные gift-card codes: encryption/KMS, reveal, повторная отправка, audit, retention и deletion?
    Почему важно: hash-only storage не позволяет доставить или повторно показать конечный код.

31. Нужно определить формат импорта manual stock: CSV columns, paste list, кодировка, типы, validation errors, dedup, rejected export.
    Почему важно: ТЗ описывает trim/dedup/result, но не точный контракт импорта.

32. Что означает `voided` для stock и когда код может быть аннулирован?
    Почему важно: inventory lifecycle и reconciliation неполные.

33. Видит ли B2C-пользователь поставщика при redeem, или supplier selection всегда внутренний механизм Admin/Ops?
    Почему важно: раскрытие supplier может влиять на UX, disputes и compliance.

34. Founder diagrams про `G-engine`, `AAAgift` и "Рога и Копыта" являются MVP scope, архитектурным вариантом или background?
    Почему важно: эти диаграммы могут расширить модель поставщиков и AML-потоки за пределы PRD/TZ.

### C. Ваучер, заказ, доставка, email и B2C support

35. Какой формат кода ваучера должен видеть пользователь и партнер: длина, группы символов, разделители, регистр, checksum, пример отображения?
    Почему важно: Rewarble reference упоминает 16-digit code, но наш формат не утвержден.

36. Нужно развести термины кодов: voucher/Redeem.io code, final supplier digital code, code shown in email/delivery, code exported to seller.
    Почему важно: founder diagrams смешивают `Digital Code`, `AAAcode` и voucher.

37. Нужно определить единый canonical status/transition map для Voucher, RedemptionOrder, FulfillmentItem, DeliveryToken и user-visible labels.
    Почему важно: в founder docs есть фрагменты (`PENDING_FULFILLMENT`, `PENDING_EXTENDED`, `MANUAL_ESCALATION`, `DELIVERED`, `FAILED`), но нет единой карты.

38. Когда пользователь видит `Failed`, когда остается `Pending`, а когда заказ уходит в `MANUAL_ESCALATION`?
    Почему важно: коммуникация и support actions различаются.

39. Нужно устранить конфликт SLA/copy: PRD/NFR требует P95 <= 30 секунд и max 2 минуты, но TZ email copy говорит "обычно 5-15 минут".
    Почему важно: acceptance, emails и marketplace approval должны обещать одно и то же.

40. Нужно определить DeliveryToken TTL, одноразовость, повторное открытие, rotation, recovery после истечения и поведение старых email-ссылок.
    Почему важно: `DeliveryToken TTL` оставлен плейсхолдером.

41. Как пользователь восстанавливает доступ к delivered order после истечения DeliveryToken?
    Почему важно: без recovery возрастет нагрузка на поддержку и риск раскрытия кодов.

42. Нужно ли подтверждение email до доставки, и принимает ли продукт риск опечатки?
    Почему важно: текущий flow описывает syntax/MX validation и отправку email, но не подтверждение владения адресом.

43. Может ли пользователь изменить email после start/confirm redeem?
    Почему важно: доставка, privacy и support search зависят от правила изменения адреса.

44. Нужно устранить конфликт email sequence: PRD говорит, что email сразу после подтверждения содержит код(ы), а TZ разделяет `Order received` без кода и `Delivered` с кодом.
    Почему важно: email templates, security и fulfillment timing зависят от финальной последовательности.

45. Нужно определить resend email limit: количество, окно времени, блокировки и user-facing copy.
    Почему важно: защита от abuse/email spam и UX должны быть тестируемыми.

46. Что происходит при hard bounce, soft bounce, spam complaint, suppression, typo или доставке неверному получателю?
    Почему важно: delivery finality, privacy incident handling и support escalation зависят от email failure policy.

47. Нужно описать безопасный текст ошибки для invalid, expired, already used voucher и technical error.
    Почему важно: ошибки не должны облегчать перебор валидных кодов.

48. Может ли пользователь вернуться к delivered order, снова введя voucher/email?
    Почему важно: модель token recovery и repeated access не определена.

49. Нужно определить rate limit/polling behavior для кнопки "Проверить статус" на `/delivery/:token`.
    Почему важно: статусные polling-запросы могут создавать нагрузку и provider/email spam.

50. Что происходит, если один voucher одновременно вводят в двух браузерах/email до подтверждения или пока первый заказ в `PENDING_FULFILLMENT`?
    Почему важно: race behavior нужен для блокировок и антифрода.

51. Нужно описать user-facing коммуникацию при `MANUAL_ESCALATION`: что видит пользователь после 6 часов, как часто получает обновления и кто отвечает за первый ответ.
    Почему важно: support SLA и доверие пользователя зависят от текста и процесса.

52. Кто является первой линией B2C support и какой фактический support email/form/chat используется?
    Почему важно: ТЗ дает routing, но `[SUPPORT EMAIL]` и операционный канал не заполнены.

53. Публичный redeem flow должен быть под единым брендом Redeem.io, partner-branded или white-label на нескольких доменах?
    Почему важно: домены, email sender, marketplace instructions, screenshots и trust copy различаются.

54. Должен ли B2C-пользователь явно принять Terms, Privacy, no refund/finality перед confirm redeem, и как хранится consent/version?
    Почему важно: legal evidence и PoD могут требовать timestamp, version, IP/device context.

55. Нужны ли B2C geo/sanctions/KYC checks, или продукт принципиально собирает только email конечника?
    Почему важно: `analytics.md` предлагает KYC triggers, а `tz.md` требует data minimization и отсутствие документов конечника.

56. Разрешено ли несовершеннолетним погашать ваучеры, и какая возрастная политика применяется?
    Почему важно: terms acceptance, privacy children и support handling требуют юридического решения.

### D. B2B кабинет, API, Admin и операции

57. Предоставьте финальный KYB-процесс/провайдера, список документов, allowed/restricted jurisdictions, UBO threshold, periodic review и legal requirements.
    Почему важно: форма KYB есть, но provider и правила проверки не утверждены.

58. Что происходит после KYB `REJECTED`: финальный отказ, повторная подача, appeal через support, срок ожидания?
    Почему важно: transitions и notifications не определены.

59. Обязательна ли 2FA для партнеров и внутренних ролей, на каких действиях, какие методы поддерживаются и как восстанавливается доступ?
    Почему важно: PRD говорит "дать возможность", но security acceptance требует точную политику.

60. Может ли партнерская организация иметь несколько пользователей и роли внутри org?
    Почему важно: ownership API keys, audit attribution и account settings зависят от multi-user модели.

61. Какой auth provider/session policy ожидается для partner cabinet и admin: passwords, SSO, TOTP, TTL, recovery, invitation, password policy?
    Почему важно: auth/session design и compliance audit должны быть утверждены.

62. Какие batch limits: min/max count, max nominal, daily velocity, per-partner approvals, manual review thresholds?
    Почему важно: fraud controls и UI validation требуют значений.

63. Может ли один VoucherBatch смешивать несколько номиналов, валют или SKU restrictions, или batch всегда один номинал и одна валюта?
    Почему важно: CSV export, invoices и PoD проще/сложнее в зависимости от batch shape.

64. Кто может создавать batch: partner self-service, Admin/Ops on behalf, или оба сценария?
    Почему важно: approval, audit и notifications отличаются.

65. Нужно определить CSV export для batch: columns, order, encoding, filename, повторная выгрузка, password/encrypted archive, download audit.
    Почему важно: партнер загружает эти коды на маркетплейс, а повторный доступ к raw codes чувствителен.

66. Какие partner notifications нужны: KYB status, payment confirmed/failed, batch ready, export downloaded, PoD request, dispute, suspend?
    Почему важно: без notifications процессы зависнут в ручной коммуникации.

67. Что включает transaction history и partner analytics: deposits, batches, redemptions, supplier costs, refunds, PoD requests, individual rows или aggregates only?
    Почему важно: founder docs требуют агрегаты без PII, но партнеру может понадобиться детализация для споров.

68. Нужно определить критерии Tier 1+ и допуск к API: volume, KYB status, manual approval, limits, contract terms.
    Почему важно: Partner API ограничен Tier 1+, но criteria не заданы.

69. Нужно согласовать конфликт Partner API scope: в PRD P7 стоит P2/Scale, а в TZ `API Issue` помечен MUST.
    Почему важно: MVP Pilot scope меняется, если API Issue требуется сразу.

70. Какой стиль public/internal API нужен: REST, GraphQL или mixed, versioning, generated clients, backward compatibility?
    Почему важно: партнерское API, admin API, webhooks и clients требуют стабильного контрактного стиля.

71. Какие webhook guarantees нужны: signature, retry schedule, event ordering, idempotency, dead-letter, manual replay?
    Почему важно: signed status webhooks упомянуты, но delivery semantics не определены.

72. Интеграции с маркетплейсами в MVP только manual CSV/listing upload, или нужен provision/reserve API flow для Eneba/Kinguin?
    Почему важно: compliance docs указывают provision API как fallback, если universal external redeem не пройдет.

73. Что именно входит в marketplace voucher package: code, redeem URL, instructions, title, brand, region, expiry, support URL, template per marketplace?
    Почему важно: partner listing и buyer delivery должны быть согласованы до batch export.

74. Нужен ли marketplace order context в redeem/support: marketplace order id, seller id, listing URL, case id?
    Почему важно: PoD и support не смогут связать redemption с marketplace dispute без контекста.

75. Что происходит при `suspend` партнера с активными unused vouchers или pending redemptions?
    Почему важно: блокировка партнера может затронуть покупателей, которые уже купили voucher.

76. Какие значения по умолчанию и единицы измерения у `daily_limit_usd` и `velocity_limit`?
    Почему важно: admin limits есть в PRD/TZ, но значения и единицы не утверждены.

77. Как отключение SKU влияет на existing reservations, pending orders, partner catalog и marketplace listing availability?
    Почему важно: PRD говорит "скрывает мгновенно", но не описывает активные заказы.

78. Что именно означает `reveal once` для Ops/Admin: кому доступно, сколько раз, какие approvals, что логируется, можно ли повторить?
    Почему важно: raw code visibility - чувствительная операция.

79. Какие дополнительные внутренние роли нужны кроме Admin/Ops/Support: Compliance, Finance, Read-only, Auditor?
    Почему важно: минимальные роли закрыты, но compliance/finance workflows требуют прав.

80. Какая ticket system используется для auto-ticket после 6 часов, или MVP делает notes/tags/in_progress/closed в admin без внешней системы?
    Почему важно: PRD говорит auto-ticket, TZ допускает минимальный ticketing.

81. Кто получает Telegram alerts, какие escalation rules, response SLA и что делать при отсутствии подтверждения?
    Почему важно: alert thresholds без owner/on-call process не работают.

82. Как конфигурируется fallback chain: global, supplier, SKU, region, partner, marketplace или mixed?
    Почему важно: PRD говорит "настраивается в Admin", но scope не задан.

83. Сколько retry attempts по одному supplier до switch supplier, и как защищаемся от duplicate issuance?
    Почему важно: retries без idempotency могут дважды списать supplier stock.

84. Что делать, если все suppliers не сработали и нет equivalent SKU?
    Почему важно: пользовательское право на получение, support SLA и refund/no-refund policy должны быть конечными.

85. Требуется ли явное согласие пользователя на альтернативный SKU, как долго система ждет ответ и что если пользователь игнорирует/отклоняет альтернативу?
    Почему важно: OOS alternative описана, но timeout/decline path не определены.

86. Какая admin surface отвечает за supplier config и fallback lifecycle: CRUD, credentials, health check, priorities, assignment by SKU/region, RBAC?
    Почему важно: нужно понимать, что именно строим в Admin для F2/F3.

87. Кто может проверять PoD и `voucher_hash`: partner, marketplace, support, только internal team; нужен ли signed/versioned export и download audit?
    Почему важно: PoD может быть непроверяемым или юридически слабым без verifier model.

### E. Compliance, fraud, безопасность и данные

88. Какие founder-доки считать нормативными, если диаграммы используют "balance" и crypto payout-потоки, а PRD/TZ запрещают wallet/balance/payout?
    Почему важно: `reedem.io.md`, `core-flow.md`, `seller-flow.md`, `analytics.md` местами конфликтуют с MUST-принципами PRD/TZ.

89. Нужно определить regulatory classification продукта по юрисдикциям: gift-card distribution, stored value, prepaid access, payment service или marketplace infrastructure.
    Почему важно: модель учета, KYB/KYT, terms и допустимые формулировки зависят от классификации.

90. Нужно выбрать server jurisdiction, legal jurisdiction и legal entity; HK упомянут, но требует подтверждения юриста.
    Почему важно: PRD оставляет EU/HK/SG и HK legal confirmation как blocker.

91. Какие страны, маркетплейсы, категории продуктов, типы партнеров и B2C jurisdictions запрещены?
    Почему важно: onboarding, catalog filtering и sanctions policy требуют allow/deny lists.

92. Какие AML/KYB/KYT требования применяются: sanctions, PEP, adverse media, UBO threshold, ongoing monitoring, provider choices и evidence retention?
    Почему важно: базовый пакет есть, но операционные правила не утверждены.

93. Нужна ли отдельная B2C KYC/light-KYC/full-KYC матрица по сумме, продукту, стране и risk factors?
    Почему важно: это новый конфликт между analytics risk advice и TZ data minimization.

94. Нужно утвердить captcha provider и точные thresholds: failed attempts N, velocity, high-risk SKU, accessibility, fallback.
    Почему важно: founder docs оставляют N и пороги пустыми.

95. Что такое "аномальная скорость операций" и какие события создают fraud hold/manual review?
    Почему важно: velocity, high-risk SKU, IP/device и supplier anomalies должны превращаться в состояния.

96. Разрешен ли device fingerprinting в выбранной юрисдикции и какая cookie/consent policy нужна?
    Почему важно: антифрод может конфликтовать с privacy.

97. Нужна blacklist/denylist model: email, IP, subnet, device, ASN, requisites, partner, supplier; retention, false positives, unblocking, access roles.
    Почему важно: analytics описывает blacklist, но нет политики.

98. Как email адреса маскируются, ищутся, хранятся, анонимизируются и используются в support search?
    Почему важно: privacy, support tooling и deletion jobs зависят от правил.

99. Какая retention policy для raw gift-card codes, voucher hashes/salts, PoD, audit logs, KYB files, supplier tx IDs и anonymized operational logs?
    Почему важно: PRD/TZ дают частичные сроки, но не полную матрицу данных.

100. Какие поля никогда не попадают в logs/PoD/analytics/vendor payloads: raw codes, secrets, full email, tokens, API keys, supplier credentials?
     Почему важно: observability и exports могут стать источником утечки.

101. Нужно определить audit log integrity model: append-only, tamper-evident, export, signing, privileged action controls, retention.
     Почему важно: incident investigation и compliance evidence требуют доверенного audit trail.

102. Как data subject rights сочетаются с disputes, legal hold, KYB files, PoD retention и anonymization?
     Почему важно: privacy requests не должны разрушить dispute evidence без legal basis.

103. Предоставьте vendor/subprocessor map: email, KYT, supplier, captcha, analytics, alerting, hosting; какие данные уходят каждому.
     Почему важно: DTO, masking, DPA и privacy notice зависят от внешних передач данных.

104. Нужны security/privacy incident playbooks: breach notification, code compromise, account/API key compromise, law-enforcement requests.
     Почему важно: raw codes, PII, API keys и audit logs требуют заранее согласованных owners и timelines.

105. Какой playbook для компрометации voucher codes до redeem после CSV export партнеру?
     Почему важно: нужно решить revoke/reissue, buyer notification, marketplace evidence и partner liability.

106. Какая модель abuse-control кроме brute-force: WAF/DDoS, global rate limits, emergency switches, provider/email/captcha cost caps, alert storms?
     Почему важно: публичный redeem может искусственно увеличивать расходы.

107. Какие RPO/RTO/SLA для value-bearing records: backup, restore, encryption key recovery, restore drills?
     Почему важно: vouchers, inventory, PoD, deposits и encrypted code storage нельзя терять без процесса.

108. Какие marketplaces явно разрешают resale/external redeem, и какое подтверждение авторизации/прав поставщика должны иметь партнеры?
     Почему важно: compliance checks для Eneba/Kinguin остаются pass with limitations.

109. Что происходит при marketplace takedown или бане аккаунта партнера?
     Почему важно: нужно решить судьбу unused vouchers, pending redemptions, buyer communication и support scripts.

110. Кто является merchant/seller of record для B2C продажи на marketplace?
     Почему важно: taxes, refunds, consumer law, support responsibility и chargeback handling зависят от seller of record.

111. Какая data-protection role model между Redeem.io, партнерами, marketplaces и vendors: controller, processor, joint-controller?
     Почему важно: DPA, privacy notice, DSAR и breach routing зависят от ролей.

### F. Метрики, monitoring, UX и контент

112. Нужно определить целевые значения: Redemption Success Rate, Failed Redemption Rate, Redeem-to-Delivery Conversion, Support Ticket Rate, Partner Retention, Gross Margin.
     Почему важно: founder docs оставляют часть метрик `[ЗАПОЛНИТЬ]`.

113. Какую peak load/throughput должна поддерживать MVP-платформа: concurrent redeems, batch size, API calls, supplier calls?
     Почему важно: queues, DB sizing и load tests требуют объема.

114. Какой допустимый процент заказов `Pending > 6h` и manually resolved?
     Почему важно: ops dashboard и alert thresholds требуют порогов.

115. При какой доле ошибок включаются Telegram alerts: global, supplier, SKU, marketplace, time window?
     Почему важно: alert policy должна быть достаточно точной для реализации.

116. Нужна observability-policy supplier health: API availability, latency, stock freshness, fail rate, fallback trigger, alert threshold.
     Почему важно: supplier monitoring упомянут, но не описан.

117. Какой email provider, sender domain, SPF/DKIM/DMARC, dedicated IP и deliverability targets нужны?
     Почему важно: email delivery - часть fulfillment и PoD.

118. Какая analytics event taxonomy нужна для redemption, batch issuance, supplier fulfillment, fraud, support, PoD и partner reports?
     Почему важно: dashboards и product metrics нельзя собирать консистентно без событий.

119. Нужно определить supported browsers/devices и accessibility standard для B2C, partner cabinet и admin.
     Почему важно: MVP web/mobile-responsive и captcha/error states требуют приемки.

120. Кто предоставляет и утверждает дизайн-артефакты: designer, AI-generation, hybrid; какие wireframes/UI-kit/states обязательны до frontend?
     Почему важно: PRD оставляет UI decision открытым.

121. Какие тексты/шаблоны нужны для landing, how-to-buy, FAQ, errors, transactional emails, KYB/payment/batch notifications и support replies?
     Почему важно: команде нельзя самостоятельно придумывать юридически чувствительные тексты.

122. Какие trust signals допустимы на landing/how-to-buy: logos, numbers, testimonials, legal badges, marketplace names?
     Почему важно: marketing claims могут быть compliance-risk.

123. Требуется ли локализация B2C redeem больше чем на один язык?
     Почему важно: content, email templates, currency UX и support language зависят от локализации.

124. Если Eneba/Kinguin не разрешат universal external redeem, должен ли MVP сразу поддерживать fallback "конкретный SKU через provision/API", или это отдельный pivot?
     Почему важно: compliance docs называют это альтернативой acceptance, но scope не утвержден.
