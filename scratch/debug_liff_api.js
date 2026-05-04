
const axios = require('axios');

async function checkLiffProducts() {
  try {
    const response = await axios.get('http://localhost:8080/api/v1/liff/products', {
      headers: {
        // We'll need a token, let's try to get one from the DB or just mock the call if possible
        // Actually, let's check the database again to see if images are associated with the right products
      }
    });
    console.log(JSON.stringify(response.data, null, 2));
  } catch (err) {
    console.error(err.message);
  }
}
