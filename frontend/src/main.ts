import { createApp } from 'vue';
import App from './App.vue';
import { router, vuetify } from './plugins';
import './assets/styles';
import 'mana-font/css/mana.min.css';

const app = createApp(App);

app.use(vuetify);
app.use(router);

app.mount('#app');
