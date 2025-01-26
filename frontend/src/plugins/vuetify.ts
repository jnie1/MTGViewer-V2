import { createVuetify, type ThemeDefinition } from 'vuetify';

const customTheme: ThemeDefinition = {
  dark: true,
  colors: {
    background: '#111519',
    'background-variant': '#212935',
    'on-background': '#e0eaea',
    'on-background-variant': '#f8f8f8',
    surface: '#232e3f',
    'on-surface': '#e0eaea',
    'on-surface-variant': '#f8f8f8',
    primary: '#fa4125',
    'primary-variant': '#fdc200',
    secondary: '#7aa5f7',
    'secondary-variant': '#02446e',
    forest: '#c4d3ca',
    island: '#b3ceea',
    plains: '#f8e7b9',
    mountain: '#eb9f82',
    swamp: '#a69f9d'
  }
};


const vuetify = createVuetify({
  theme: {
    defaultTheme: 'customTheme',
    themes: { customTheme },
  },
});

export default vuetify;
