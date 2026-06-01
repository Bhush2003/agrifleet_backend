package models

// Permission constants for granular RBAC
type Permission string

const (
	// Employee
	PermEmployeeCreate Permission = "employee:create"
	PermEmployeeUpdate Permission = "employee:update"
	PermEmployeeDelete Permission = "employee:delete"
	PermEmployeeView   Permission = "employee:view"
	PermEmployeeExport Permission = "employee:export"

	// Attendance
	PermAttendanceMark     Permission = "attendance:mark"
	PermAttendanceApprove  Permission = "attendance:approve"
	PermAttendanceModify   Permission = "attendance:modify"
	PermAttendanceLock     Permission = "attendance:lock"

	// Payroll
	PermPayrollView     Permission = "payroll:view"
	PermPayrollGenerate Permission = "payroll:generate"
	PermPayrollFinalize Permission = "payroll:finalize"
	PermPayrollLock     Permission = "payroll:lock"

	// Vehicle
	PermVehicleCreate Permission = "vehicle:create"
	PermVehicleUpdate Permission = "vehicle:update"
	PermVehicleDelete Permission = "vehicle:delete"
	PermVehicleLock   Permission = "vehicle:lock"

	// Diesel
	PermDieselCreate   Permission = "diesel:create"
	PermDieselApprove  Permission = "diesel:approve"
	PermDieselOverride Permission = "diesel:override"

	// Expense
	PermExpenseCreate  Permission = "expense:create"
	PermExpenseApprove Permission = "expense:approve"
	PermExpenseAudit   Permission = "expense:audit"

	// Income
	PermIncomeCreate  Permission = "income:create"
	PermIncomeApprove Permission = "income:approve"

	// Reports
	PermReportView       Permission = "report:view"
	PermReportExport     Permission = "report:export"
	PermReportFinancial  Permission = "report:financial"
	PermReportAudit      Permission = "report:audit"

	// Settings
	PermSettingsManage Permission = "settings:manage"

	// Backup
	PermBackupCreate  Permission = "backup:create"
	PermBackupRestore Permission = "backup:restore"

	// Audit
	PermAuditView Permission = "audit:view"

	// Approval
	PermApprovalManage Permission = "approval:manage"

	// Month Lock
	PermMonthLock Permission = "month:lock"

	// Soft Delete / Restore
	PermRecycleBin Permission = "recycle:manage"
)

// RolePermissions maps each role to its allowed permissions
var RolePermissions = map[UserRole][]Permission{
	RoleAdmin: {
		PermEmployeeCreate, PermEmployeeUpdate, PermEmployeeDelete,
		PermEmployeeView, PermEmployeeExport,
		PermAttendanceMark, PermAttendanceApprove, PermAttendanceModify, PermAttendanceLock,
		PermPayrollView, PermPayrollGenerate, PermPayrollFinalize, PermPayrollLock,
		PermVehicleCreate, PermVehicleUpdate, PermVehicleDelete, PermVehicleLock,
		PermDieselCreate, PermDieselApprove, PermDieselOverride,
		PermExpenseCreate, PermExpenseApprove, PermExpenseAudit,
		PermIncomeCreate, PermIncomeApprove,
		PermReportView, PermReportExport, PermReportFinancial, PermReportAudit,
		PermSettingsManage,
		PermBackupCreate, PermBackupRestore,
		PermAuditView,
		PermApprovalManage,
		PermMonthLock,
		PermRecycleBin,
	},
	RoleManager: {
		PermEmployeeCreate, PermEmployeeUpdate, PermEmployeeView,
		PermAttendanceMark,
		PermPayrollView,
		PermVehicleCreate, PermVehicleUpdate,
		PermDieselCreate,
		PermExpenseCreate,
		PermIncomeCreate,
		PermReportView, PermReportExport,
	},
	RoleSupervisor: {
		PermEmployeeView,
		PermAttendanceMark,
		PermPayrollView,
		PermVehicleUpdate,
		PermDieselCreate,
		PermExpenseCreate,
		PermReportView,
	},
	RoleDriver: {
		PermReportView,
	},
	RoleAccountant: {
		PermPayrollView, PermPayrollGenerate,
		PermExpenseCreate, PermExpenseApprove,
		PermIncomeCreate,
		PermReportView, PermReportExport, PermReportFinancial,
	},
}

// HasPermission checks if a role has a specific permission
func HasPermission(role UserRole, perm Permission) bool {
	perms, ok := RolePermissions[role]
	if !ok {
		return false
	}
	for _, p := range perms {
		if p == perm {
			return true
		}
	}
	return false
}
