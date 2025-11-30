export const formatThaiDateTime = (isoString) => {
  if (!isoString) return "-";
  const date = new Date(isoString);

  return date.toLocaleString("th-TH", {
    year: "numeric",
    month: "short",   // พ.ย.
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
};

export const formatThaiDate = (isoString) => {
  if (!isoString) return "-";
  const date = new Date(isoString);

  return date.toLocaleDateString("th-TH", {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
};

export const formatThaiTime = (isoString) => {
  if (!isoString) return "-";
  const date = new Date(isoString);

  return date.toLocaleTimeString("th-TH", {
    hour: "2-digit",
    minute: "2-digit",
  });
};
