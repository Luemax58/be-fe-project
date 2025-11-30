import TenantNavbar from "../components/TenantNavbar";

export default function TenantLayout({ children }) {
  return (
    <div className="min-h-screen bg-gray-50">
      <TenantNavbar />
      <div className="p-6">
        {children}
      </div>
    </div>
  );
}
