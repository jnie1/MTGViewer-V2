import { createVuetify, type ThemeDefinition } from 'vuetify';
import * as components from 'vuetify/components';
import * as directives from 'vuetify/directives';

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

const vuetify = createVuetify({
  components,
  directives,
  theme: {
    defaultTheme: 'customTheme',
    themes: { customTheme },
  },
});

export default vuetify;
