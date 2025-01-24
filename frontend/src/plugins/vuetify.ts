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
  },
};

const manaColorTheme: ThemeDefinition = {
  dark: true,
  colors: {
    forest: '#C4D3CA',
    island: '#B3CEEA',
    plains: '#F8E7B9',
    mountain: '#EB9F82',
    swamp: '#A69F9D'
  }
}

const vuetify = createVuetify({
  theme: {
    defaultTheme: 'customTheme',
    themes: { customTheme, manaColorTheme },
  },
});

export default vuetify;
