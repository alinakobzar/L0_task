const input = document.getElementById('orderId');
const button = document.getElementById('buttonId');
const result = document.getElementById('result');

button.addEventListener('click', () => {
    const orderId = input.value.trim();
    if (!orderId) {
        result.textContent = 'ID заказа не может быть пустым';
        return;
    }
    fetchOrder(orderId);
});

async function fetchOrder(id) {
    result.textContent = 'Загрузка...';
    try {
        const res = await fetch(`http://localhost:8080/orders/${id}`);
        if (!res.ok) {
            result.textContent = `Ошибка: ${res.status}`;
            return;
        }
        const data = await res.json();
        result.textContent = JSON.stringify(data, null, 2);
    } catch (err) {
        result.textContent = 'Ошибка сети: ' + err.message;
    }
}