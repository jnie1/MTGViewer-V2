export const capitalize = (str: string | null | undefined) => {
  if (!str) return '';
  return str[0].toUpperCase() + str.slice(1);
};
