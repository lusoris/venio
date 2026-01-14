"use client";

import { useAuth } from "@/contexts/AuthContext";
import { redirect } from "next/navigation";
import { useState, useEffect } from "react";
import RoleManagementTable from "@/components/admin/RoleManagementTable";
import RoleFormModal from "@/components/admin/RoleFormModal";

interface Role {
  id: string;
  name: string;
  description: string;
  permissions: string[];
  user_count: number;
  created_at: string;
}

export default function RolesPage() {
  const { user, isAuthenticated, token } = useAuth();
  const [roles, setRoles] = useState<Role[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showForm, setShowForm] = useState(false);

  if (!isAuthenticated || !user) {
    redirect("/login");
  }

  if (user.roles && !user.roles.includes("admin")) {
    redirect("/dashboard");
  }

  useEffect(() => {
    fetchRoles();
  }, [token]);

  const fetchRoles = async () => {
    try {
      setLoading(true);
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/api/v1/admin/roles`,
        {
          headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "application/json",
          },
        }
      );

      if (!response.ok) {
        throw new Error(`Failed to fetch roles: ${response.statusText}`);
      }

      const data = await response.json();
      setRoles(data.roles || []);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch roles");
    } finally {
      setLoading(false);
    }
  };

  const handleRoleUpdated = () => {
    fetchRoles();
    setShowForm(false);
  };

  return (
    <div>
      <div className="mb-6 flex justify-between items-start">
        <div>
          <h1 className="text-3xl font-bold text-white">Role Management</h1>
          <p className="text-gray-400 mt-2">
            Create and manage roles with specific permissions
          </p>
        </div>
        <button
          onClick={() => setShowForm(true)}
          className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg"
        >
          Create Role
        </button>
      </div>

      {error && (
        <div className="bg-red-900 border border-red-700 text-red-100 px-4 py-3 rounded-lg mb-4">
          {error}
        </div>
      )}

      {loading ? (
        <div className="bg-gray-800 rounded-lg p-8 text-center border border-gray-700">
          <p className="text-gray-400">Loading roles...</p>
        </div>
      ) : (
        <RoleManagementTable
          roles={roles}
          onRefresh={fetchRoles}
          onRoleUpdated={handleRoleUpdated}
        />
      )}

      {showForm && (
        <RoleFormModal
          onClose={() => setShowForm(false)}
          onSuccess={handleRoleUpdated}
        />
      )}
    </div>
  );
}
