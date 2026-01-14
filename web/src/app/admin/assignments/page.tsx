"use client";

import { useAuth } from "@/contexts/AuthContext";
import { redirect } from "next/navigation";
import { useState, useEffect } from "react";
import AssignmentManagementTable from "@/components/admin/AssignmentManagementTable";

interface UserRole {
  id: string;
  user_id: string;
  user_email: string;
  role_id: string;
  role_name: string;
  assigned_at: string;
  assigned_by: string;
}

export default function AssignmentsPage() {
  const { user, isAuthenticated, token } = useAuth();
  const [assignments, setAssignments] = useState<UserRole[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  if (!isAuthenticated || !user) {
    redirect("/login");
  }

  if (user.roles && !user.roles.includes("admin")) {
    redirect("/dashboard");
  }

  useEffect(() => {
    fetchAssignments();
  }, [token]);

  const fetchAssignments = async () => {
    try {
      setLoading(true);
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/api/v1/admin/user-roles`,
        {
          headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "application/json",
          },
        }
      );

      if (!response.ok) {
        throw new Error(`Failed to fetch assignments: ${response.statusText}`);
      }

      const data = await response.json();
      setAssignments(data.assignments || []);
      setError(null);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Failed to fetch assignments"
      );
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-white">Role Assignments</h1>
        <p className="text-gray-400 mt-2">
          Manage user-to-role assignments and permissions
        </p>
      </div>

      {error && (
        <div className="bg-red-900 border border-red-700 text-red-100 px-4 py-3 rounded-lg mb-4">
          {error}
        </div>
      )}

      {loading ? (
        <div className="bg-gray-800 rounded-lg p-8 text-center border border-gray-700">
          <p className="text-gray-400">Loading assignments...</p>
        </div>
      ) : (
        <AssignmentManagementTable
          assignments={assignments}
          onRefresh={fetchAssignments}
        />
      )}
    </div>
  );
}
