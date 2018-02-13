tok=$(printf "$EMAIL:$TOKEN" | base64 -w 0)
echo $tok
curl -X POST "https://selly.gg/api/v2/pay" \
  -H "Content-Type: application/json" \
  -H "Authorization: Basic $tok" \
  -d '{"title":"Selly Pay Example", "gateway":"Bitcoin", "email":"customer@email.com", "value":"10", "currency":"USD", "return_url":"https://website.com/return", "webhook_url":"https://website.com/webhook?secret=cEZMeEVlTz"}'
