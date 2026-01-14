"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useAuth } from "@/contexts/AuthContext";
import { useRouter } from "next/navigation";

export default function AdminSidebar() {
  const pathname = usePathname();
  const { logout } = useAuth();
  const router = useRouter();

  const handleLogout = () => {
    logout();
    router.push("/login");
  };

  const navItems = [
    {
      name: "Dashboard",
      href: "/admin",
      icon: "📊",
    },
    {
      name: "Users",
      href: "/admin/users",
      icon: "👥",
    },
    {
      name: "Roles",
      href: "/admin/roles",
      icon: "🔐",
    },
    {
      name: "Permissions",
      href: "/admin/permissions",
      icon: "📋",
    },
    {
      name: "Assignments",
      href: "/admin/assignments",
      icon: "🔗",
    },
  ];

  return (
    <aside className="w-64 bg-gray-800 border-r border-gray-700 flex flex-col">
      {/* Logo */}
      <div className="p-6 border-b border-gray-700">
        <h2 className="text-xl font-bold text-white">Venio Admin</h2>
        <p className="text-gray-400 text-xs mt-1">Management Panel</p>
      </div>

      {/* Navigation */}
      <nav className="flex-1 p-4 space-y-2">
        {navItems.map((item) => {
          const isActive = pathname === item.href;
          return (
            <Link
              key={item.href}
              href={item.href}
              className={`flex items-center px-4 py-2 rounded-lg transition ${
                isActive
                  ? "bg-blue-600 text-white"
                  : "text-gray-400 hover:text-white hover:bg-gray-700"
              }`}
            >
              <span className="mr-3">{item.icon}</span>
              {item.name}
            </Link>
          );
        })}
      </nav>

      {/* Footer */}
      <div className="p-4 border-t border-gray-700 space-y-2">
        <Link
          href="/dashboard"
          className="flex items-center px-4 py-2 text-gray-400 hover:text-white hover:bg-gray-700 rounded-lg transition"
        >
          <span className="mr-3">←</span>
          Back to App
        </Link>
        <button
          onClick={handleLogout}
          className="w-full flex items-center px-4 py-2 text-gray-400 hover:text-white hover:bg-gray-700 rounded-lg transition"
        >
          <span className="mr-3">🚪</span>
          Logout
        </button>
      </div>
    </aside>
  );
}
