/**
 * 根据月份获取当前季节
 * @returns {string} spring | summer | autumn | winter
 */
export const getCurrentSeason = () => {
  const month = new Date().getMonth() + 1; // getMonth() 返回 0-11
  
  if (month >= 3 && month <= 5) return 'spring';
  if (month >= 6 && month <= 8) return 'summer';
  if (month >= 9 && month <= 11) return 'autumn';
  return 'winter';
};