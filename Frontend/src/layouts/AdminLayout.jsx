import Navbar from "../components/Navbar";

export default function AdminLayout({ children }) {
  return (
    <div>
      <Navbar role="owner" />
      <div className="p-6">{children}</div>
    </div>
  );
}
