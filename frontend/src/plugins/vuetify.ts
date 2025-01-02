import { createVuetify, type ThemeDefinition } from 'vuetify';
import * as components from 'vuetify/components';
import * as directives from 'vuetify/directives';

const customTheme: ThemeDefinition = {
  dark: true,
  colors: {
    background: '#111519',
    surface: '#232e3f',
    primary: '#fa4125',
    secondary: '#7aa5f7',
  }
};

const vuetify = createVuetify({
  components,
  directives,
  theme: {
    defaultTheme: 'customTheme',
    themes: { customTheme }
  },
});

export default vuetify;
