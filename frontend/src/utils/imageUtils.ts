// Массив путей к ассетам. Сервис для генерации картинок работает только с VPN, поэтому во избежании ошибок вручную добавили картинки.
const images = [
  '/images/products/1.jpg',
  '/images/products/2.jpg',
  '/images/products/3.jpg',
  '/images/products/4.jpg',
  '/images/products/5.jpg',
  '/images/products/6.jpg',
  '/images/products/7.jpg',
  '/images/products/8.jpg',
  '/images/products/9.jpg',
  '/images/products/10.jpg',
  '/images/products/11.jpg',
  '/images/products/12.jpg',
  '/images/products/13.jpg',
  '/images/products/14.jpg',
  '/images/products/15.jpg',
  '/images/products/16.jpg',
  '/images/products/17.jpg',
  '/images/products/18.jpg',
  '/images/products/19.jpg',
  '/images/products/20.jpg',
];

export const getProductImage = (seed: string): string => {
  let hash = 0;
  for (let i = 0; i < seed.length; i++) {
    hash = (hash << 5) - hash + seed.charCodeAt(i);
    hash |= 0;
  }
  const index = Math.abs(hash) % images.length;
  return images[index];
};
export const getRandomProductImage = (): string => {
  const index = Math.floor(Math.random() * images.length);
  return images[index];
};