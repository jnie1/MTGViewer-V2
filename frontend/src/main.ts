import { createApp } from 'vue';
import App from './App.vue';
import { router, vuetify } from './plugins';
import './assets/styles';

const app = createApp(App);

app.use(vuetify);
app.use(router);

app.mount('#app');
